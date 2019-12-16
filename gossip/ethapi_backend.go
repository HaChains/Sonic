package gossip

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/bloombits"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	notify "github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	errors2 "github.com/pkg/errors"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/gossip/gasprice"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis/sfc"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis/sfc/sfcpos"
	"github.com/Fantom-foundation/go-lachesis/tracing"
)

var ErrNotImplemented = func(name string) error { return errors.New(name + " method is not implemented yet") }

// EthAPIBackend implements ethapi.Backend.
type EthAPIBackend struct {
	extRPCEnabled bool
	svc           *Service
	state         *EvmStateReader
	gpo           *gasprice.Oracle
	mux           *notify.TypeMux
}

// ChainConfig returns the active chain configuration.
func (b *EthAPIBackend) ChainConfig() *params.ChainConfig {
	return params.AllEthashProtocolChanges
}

func (b *EthAPIBackend) CurrentBlock() *evmcore.EvmBlock {
	return b.state.CurrentBlock()
}

func (b *EthAPIBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*evmcore.EvmHeader, error) {
	blk, err := b.BlockByNumber(ctx, number)
	return blk.Header(), err
}

func (b *EthAPIBackend) HeaderByHash(ctx context.Context, h common.Hash) (*evmcore.EvmHeader, error) {
	index := b.svc.store.GetBlockIndex(hash.Event(h))
	if index == nil {
		return nil, nil
	}
	return b.HeaderByNumber(ctx, rpc.BlockNumber(*index))
}

func (b *EthAPIBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*evmcore.EvmBlock, error) {
	if number == rpc.PendingBlockNumber {
		return nil, errors.New("pending block request isn't allowed")
	}
	// Otherwise resolve and return the block
	var blk *evmcore.EvmBlock
	if number == rpc.LatestBlockNumber {
		blk = b.state.CurrentBlock()
	} else {
		n := uint64(number.Int64())
		blk = b.state.GetBlock(common.Hash{}, n)
	}

	return blk, nil
}

func (b *EthAPIBackend) StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*state.StateDB, *evmcore.EvmHeader, error) {
	if number == rpc.PendingBlockNumber {
		return nil, nil, errors.New("pending block request isn't allowed")
	}
	var header *evmcore.EvmHeader
	if number == rpc.LatestBlockNumber {
		header = &b.state.CurrentBlock().EvmHeader
	} else {
		header = b.state.GetHeader(common.Hash{}, uint64(number))
	}
	if header == nil {
		return nil, nil, errors.New("header not found")
	}
	stateDb := b.svc.store.StateDB(header.Root)
	return stateDb, header, nil
}

// decodeShortEventID decodes ShortID
// example of a ShortID: "5:26:a2395846", where 5 is epoch, 26 is lamport, a2395846 are first bytes of the hash
// s is a string splitted by ":" separator
func decodeShortEventID(s []string) (idx.Epoch, idx.Lamport, []byte, error) {
	if len(s) != 3 {
		return 0, 0, nil, errors.New("incorrect format of short event ID (need Epoch:Lamport:Hash")
	}
	epoch, err := strconv.ParseUint(s[0], 10, 32)
	if err != nil {
		return 0, 0, nil, errors2.Wrap(err, "short hash parsing error (lamport)")
	}
	lamport, err := strconv.ParseUint(s[1], 10, 32)
	if err != nil {
		return 0, 0, nil, errors2.Wrap(err, "short hash parsing error (lamport)")
	}
	return idx.Epoch(epoch), idx.Lamport(lamport), common.FromHex(s[2]), nil
}

// GetFullEventID "converts" ShortID to full event's hash, by searching in events DB.
func (b *EthAPIBackend) GetFullEventID(shortEventID string) (hash.Event, error) {
	s := strings.Split(shortEventID, ":")
	if len(s) == 1 {
		// it's a full hash
		return hash.HexToEventHash(shortEventID), nil
	}
	// short hash
	epoch, lamport, prefix, err := decodeShortEventID(s)
	if err != nil {
		return hash.Event{}, err
	}

	b.svc.engineMu.RLock() // lock because of iteration
	defer b.svc.engineMu.RUnlock()

	options := b.svc.store.FindEventHashes(epoch, lamport, prefix)
	if len(options) == 0 {
		return hash.Event{}, errors.New("event not found by short ID")
	}
	if len(options) > 1 {
		return hash.Event{}, errors.New("there're multiple events with the same short ID, please use full ID")
	}
	return options[0], nil
}

// GetEvent returns Lachesis event by hash or short ID.
func (b *EthAPIBackend) GetEvent(ctx context.Context, shortEventID string) (*inter.Event, error) {
	id, err := b.GetFullEventID(shortEventID)
	if err != nil {
		return nil, err
	}
	return b.svc.store.GetEvent(id), nil
}

// GetEventHeader returns the Lachesis event header by hash or short ID.
func (b *EthAPIBackend) GetEventHeader(ctx context.Context, shortEventID string) (*inter.EventHeaderData, error) {
	id, err := b.GetFullEventID(shortEventID)
	if err != nil {
		return nil, err
	}
	epoch := id.Epoch()
	if epoch != b.svc.engine.GetEpoch() {
		return nil, errors.New("event headers are stored only for current epoch")
	}
	return b.svc.store.GetEventHeader(epoch, id), nil
}

// GetConsensusTime returns event's consensus time, if event is confirmed.
func (b *EthAPIBackend) GetConsensusTime(ctx context.Context, shortEventID string) (inter.Timestamp, error) {
	id, err := b.GetFullEventID(shortEventID)
	if err != nil {
		return 0, err
	}
	return b.svc.engine.GetConsensusTime(id)
}

// GetHeads returns IDs of all the epoch events with no descendants.
// * When epoch is -1 the heads for latest epoch are returned.
func (b *EthAPIBackend) GetHeads(ctx context.Context, epoch rpc.BlockNumber) (heads hash.Events, err error) {
	current := b.svc.engine.GetEpoch()

	var requested idx.Epoch
	switch {
	case epoch == rpc.LatestBlockNumber:
		requested = current
	case epoch >= 0 && idx.Epoch(epoch) <= current:
		requested = idx.Epoch(epoch)
	default:
		err = errors.New("epoch is not in range")
		return
	}

	if requested == current {
		heads = b.svc.store.GetHeads(requested)
	} else {
		num, ok := b.svc.store.GetPacksNum(requested)
		if !ok {
			err = errors.New("epoch is not found")
			return
		}
		packInfo := b.svc.store.GetPackInfo(requested, num-1)
		if packInfo == nil {
			err = errors.New("epoch is not found")
			return
		}
		heads = packInfo.Heads
	}

	if heads == nil {
		heads = hash.Events{}
	}

	return
}

func (b *EthAPIBackend) GetHeader(ctx context.Context, h common.Hash) *evmcore.EvmHeader {
	header, err := b.HeaderByHash(ctx, h)
	if err != nil {
		return nil
	}
	return header
}

func (b *EthAPIBackend) GetBlock(ctx context.Context, h common.Hash) (*evmcore.EvmBlock, error) {
	index := b.svc.store.GetBlockIndex(hash.Event(h))
	if index == nil {
		return nil, nil
	}
	return b.BlockByNumber(ctx, rpc.BlockNumber(*index))
}

func (b *EthAPIBackend) GetReceipts(ctx context.Context, number rpc.BlockNumber) (types.Receipts, error) {
	if !b.svc.config.TxIndex {
		return nil, errors.New("transactions index is disabled (enable TxIndex and re-process the DAG)")
	}

	if number == rpc.PendingBlockNumber {
		return nil, errors.New("pending block request isn't allowed")
	}
	if number == rpc.LatestBlockNumber {
		header := b.state.CurrentHeader()
		number = rpc.BlockNumber(header.Number.Uint64())
	}

	receipts := b.svc.store.GetReceipts(idx.Block(number))
	return receipts, nil
}

func (b *EthAPIBackend) GetLogs(ctx context.Context, number rpc.BlockNumber) ([][]*types.Log, error) {
	receipts, err := b.GetReceipts(ctx, number)
	if receipts == nil || err != nil {
		return nil, err
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *EthAPIBackend) GetTd(blockHash common.Hash) *big.Int {
	return big.NewInt(0)
}

func (b *EthAPIBackend) GetEVM(ctx context.Context, msg evmcore.Message, state *state.StateDB, header *evmcore.EvmHeader) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := evmcore.NewEVMContext(msg, header, b.state, nil)
	config := params.AllEthashProtocolChanges
	return vm.NewEVM(context, state, config, vm.Config{}), vmError, nil
}

func (b *EthAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	err := b.svc.txpool.AddLocal(signedTx)
	if err == nil {
		// NOTE: only sent txs tracing, see TxPool.addTxs() for all
		tracing.StartTx(signedTx.Hash(), "EthAPIBackend.SendTx()")
		// TODO: txLatency cleaning, possible memory leak
		txLatency.Start(signedTx.Hash())
	}
	return err
}

func (b *EthAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.svc.txpool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *EthAPIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.svc.txpool.Get(hash)
}

func (b *EthAPIBackend) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, uint64, uint64, error) {
	if !b.svc.config.TxIndex {
		return nil, 0, 0, errors.New("transactions index is disabled (enable TxIndex and re-process the DAG)")
	}

	position := b.svc.store.GetTxPosition(txHash)
	if position == nil {
		return nil, 0, 0, nil
	}

	event := b.svc.store.GetEvent(position.Event)
	if position.EventOffset > uint32(event.Transactions.Len()) {
		return nil, 0, 0, fmt.Errorf("transactions index is corrupted (offset is larger than number of txs in event), event=%s, txid=%s, block=%d, offset=%d, txs_num=%d",
			position.Event.String(),
			txHash.String(),
			position.Block,
			position.EventOffset,
			event.Transactions.Len())
	}

	tx := event.Transactions[position.EventOffset]
	return tx, uint64(position.Block), uint64(position.BlockOffset), nil
}

func (b *EthAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.svc.txpool.Nonce(addr), nil
}

func (b *EthAPIBackend) Stats() (pending int, queued int) {
	return b.svc.txpool.Stats()
}

func (b *EthAPIBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.svc.txpool.Content()
}

func (b *EthAPIBackend) SubscribeNewTxsNotify(ch chan<- evmcore.NewTxsNotify) notify.Subscription {
	return b.svc.txpool.SubscribeNewTxsNotify(ch)
}

func (b *EthAPIBackend) Progress() PeerProgress {
	return b.svc.pm.myProgress()
}

func (b *EthAPIBackend) ProtocolVersion() int {
	return int(ProtocolVersions[len(ProtocolVersions)-1])
}

func (b *EthAPIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *EthAPIBackend) ChainDb() ethdb.Database {
	return b.svc.store.table.Evm
}

func (b *EthAPIBackend) NotifyMux() *notify.TypeMux {
	return b.mux
}

func (b *EthAPIBackend) AccountManager() *accounts.Manager {
	return b.svc.AccountManager()
}

func (b *EthAPIBackend) ExtRPCEnabled() bool {
	return b.extRPCEnabled
}

func (b *EthAPIBackend) RPCGasCap() *big.Int {
	return b.svc.config.RPCGasCap
}

func (b *EthAPIBackend) BloomStatus() (uint64, uint64) {
	// TODO: implement or disable it. Origin:
	/*
		sections, _, _ := b.svc.bloomIndexer.Sections()
		return params.BloomBitsBlocks, sections
	*/
	return 0, 0
}

func (b *EthAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	// TODO: implement or disable it. Origin:
	/*
		for i := 0; i < bloomFilterThreads; i++ {
			go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.svc.bloomRequests)
		}
	*/
}

// CurrentEpoch returns current epoch number.
func (b *EthAPIBackend) CurrentEpoch(ctx context.Context) idx.Epoch {
	return b.svc.engine.GetEpoch()
}

// GetEpochStats returns epoch statistics.
// * When epoch is -1 the statistics for latest epoch is returned.
func (b *EthAPIBackend) GetEpochStats(ctx context.Context, requestedEpoch rpc.BlockNumber) (*sfctype.EpochStats, idx.Epoch, error) {
	var epoch idx.Epoch
	if requestedEpoch == rpc.PendingBlockNumber {
		epoch = pendingEpoch
	} else if requestedEpoch == rpc.LatestBlockNumber {
		epoch = b.CurrentEpoch(ctx) - 1
	} else {
		epoch = idx.Epoch(requestedEpoch)
	}
	if epoch == b.CurrentEpoch(ctx) {
		return nil, 0, errors.New("current epoch isn't sealed yet, request pending epoch")
	}

	return b.svc.store.GetEpochStats(epoch), epoch, nil
}

// GetValidationScore returns staker's ValidationScore.
func (b *EthAPIBackend) GetValidationScore(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	if !b.svc.store.HasSfcStaker(stakerID) {
		return nil, nil
	}
	return new(big.Int).SetUint64(b.svc.store.GetActiveValidationScore(stakerID)), nil
}

// GetOriginationScore returns staker's OriginationScore.
func (b *EthAPIBackend) GetOriginationScore(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	if !b.svc.store.HasSfcStaker(stakerID) {
		return nil, nil
	}
	return new(big.Int).SetUint64(b.svc.store.GetActiveOriginationScore(stakerID)), nil
}

// GetStakerPoI returns staker's PoI.
func (b *EthAPIBackend) GetStakerPoI(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	if !b.svc.store.HasSfcStaker(stakerID) {
		return nil, nil
	}
	return new(big.Int).SetUint64(b.svc.store.GetStakerPOI(stakerID)), nil
}

// GetValidatingPower returns staker's ValidatingPower.
func (b *EthAPIBackend) GetValidatingPower(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	if !b.svc.store.HasSfcStaker(stakerID) {
		return nil, nil
	}
	header := b.state.CurrentHeader()
	statedb := b.svc.store.StateDB(header.Root)

	epoch := b.svc.engine.GetEpoch()
	epochPosition := sfcpos.EpochSnapshot(epoch - 1)
	validatorPosition := epochPosition.ValidatorMerit(stakerID)
	validatingPower256 := statedb.GetState(sfc.ContractAddress, validatorPosition.ValidatingPower())

	return new(big.Int).SetBytes(validatingPower256.Bytes()), nil
}

// GetDowntime returns staker's Downtime.
func (b *EthAPIBackend) GetDowntime(ctx context.Context, stakerID idx.StakerID) (idx.Block, inter.Timestamp, error) {
	epoch := b.svc.engine.GetEpoch()
	if !b.svc.store.HasEpochValidator(epoch, stakerID) {
		return 0, 0, errors.New("staker isn't validator")
	}

	missed := b.svc.store.GetBlocksMissed(stakerID)
	return missed.Num, missed.Period, nil
}

// GetStaker returns SFC staker's info
func (b *EthAPIBackend) GetStaker(ctx context.Context, stakerID idx.StakerID) (*sfctype.SfcStaker, error) {
	staker := b.svc.store.GetSfcStaker(stakerID)
	if staker == nil {
		return nil, nil
	}
	staker.IsValidator = b.svc.engine.GetValidators().Exists(stakerID)
	return staker, nil
}

// GetStakerID returns SFC staker's Id by address
func (b *EthAPIBackend) GetStakerID(ctx context.Context, addr common.Address) (idx.StakerID, error) {
	header := b.state.CurrentHeader()
	statedb := b.svc.store.StateDB(header.Root)

	position := sfcpos.StakerID(addr)
	stakerID256 := statedb.GetState(sfc.ContractAddress, position)

	return idx.StakerID(new(big.Int).SetBytes(stakerID256.Bytes()).Uint64()), nil
}

// GetStakers returns SFC stakers info
func (b *EthAPIBackend) GetStakers(ctx context.Context) ([]sfctype.SfcStakerAndID, error) {
	b.svc.engineMu.RLock() // lock because of iteration
	defer b.svc.engineMu.RUnlock()

	stakers := make([]sfctype.SfcStakerAndID, 0, 200)
	b.svc.store.ForEachSfcStaker(func(it sfctype.SfcStakerAndID) {
		it.Staker.IsValidator = b.svc.engine.GetValidators().Exists(it.StakerID)
		stakers = append(stakers, it)
	})
	return stakers, nil
}

// GetDelegatorsOf returns SFC delegators who delegated to a staker
func (b *EthAPIBackend) GetDelegatorsOf(ctx context.Context, stakerID idx.StakerID) ([]sfctype.SfcDelegatorAndAddr, error) {
	b.svc.engineMu.RLock() // lock because of iteration
	defer b.svc.engineMu.RUnlock()

	delegators := make([]sfctype.SfcDelegatorAndAddr, 0, 200)
	// TODO add additional DB index
	b.svc.store.ForEachSfcDelegator(func(it sfctype.SfcDelegatorAndAddr) {
		if it.Delegator.ToStakerID == stakerID {
			delegators = append(delegators, it)
		}
	})
	return delegators, nil
}

// GetDelegator returns SFC delegator info
func (b *EthAPIBackend) GetDelegator(ctx context.Context, addr common.Address) (*sfctype.SfcDelegator, error) {
	return b.svc.store.GetSfcDelegator(addr), nil
}

// GetDelegatorClaimedRewards returns sum of claimed rewards in past, by this delegator
func (b *EthAPIBackend) GetDelegatorClaimedRewards(ctx context.Context, addr common.Address) (*big.Int, error) {
	return b.svc.store.GetDelegatorClaimedRewards(addr), nil
}

// GetStakerClaimedRewards returns sum of claimed rewards in past, by this staker
func (b *EthAPIBackend) GetStakerClaimedRewards(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	return b.svc.store.GetStakerClaimedRewards(stakerID), nil
}

// GetStakerDelegatorsClaimedRewards returns sum of claimed rewards in past, by this delegators of this staker
func (b *EthAPIBackend) GetStakerDelegatorsClaimedRewards(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	return b.svc.store.GetStakerDelegatorsClaimedRewards(stakerID), nil
}