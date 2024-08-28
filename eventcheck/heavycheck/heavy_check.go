package heavycheck

import (
	"bytes"
	"errors"
	"runtime"
	"sync"

	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/Fantom-foundation/go-opera/eventcheck/epochcheck"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/inter/validatorpk"
)

var (
	ErrWrongEventSig            = errors.New("event has wrong signature")
	ErrMalformedTxSig           = errors.New("tx has wrong signature")
	ErrWrongPayloadHash         = errors.New("event has wrong payload hash")
	ErrPubkeyChanged            = errors.New("validator pubkey has changed, cannot create BVs/EV for older epochs")
	ErrInvalidBridgeVote        = errors.New("invalid bridge vote")

	errTerminated = errors.New("terminated") // internal err
)

// Reader is accessed by the validator to get the current state.
type Reader interface {
	GetEpochPubKeys() (map[idx.ValidatorID]validatorpk.PubKey, idx.Epoch)
	GetEpochPubKeysOf(idx.Epoch) map[idx.ValidatorID]validatorpk.PubKey
	GetEpochBlockStart(idx.Epoch) idx.Block
}

// Checker which requires only parents list + current epoch info
type Checker struct {
	config   Config
	txSigner types.Signer
	reader   Reader

	tasksQ chan *taskData
	quit   chan struct{}
	wg     sync.WaitGroup
}

type taskData struct {
	event inter.EventPayloadI

	onValidated func(error)
}

// New validator which performs heavy checks, related to signatures validation and Merkle tree validation
func New(config Config, reader Reader, txSigner types.Signer) *Checker {
	if config.Threads == 0 {
		config.Threads = runtime.NumCPU()
		if config.Threads > 1 {
			config.Threads--
		}
		if config.Threads < 1 {
			config.Threads = 1
		}
	}
	return &Checker{
		config:   config,
		txSigner: txSigner,
		reader:   reader,
		tasksQ:   make(chan *taskData, config.MaxQueuedTasks),
		quit:     make(chan struct{}),
	}
}

func (v *Checker) Start() {
	for i := 0; i < v.config.Threads; i++ {
		v.wg.Add(1)
		go v.loop()
	}
}

func (v *Checker) Stop() {
	close(v.quit)
	v.wg.Wait()
}

func (v *Checker) Overloaded() bool {
	return len(v.tasksQ) > v.config.MaxQueuedTasks/2
}

func (v *Checker) EnqueueEvent(e inter.EventPayloadI, onValidated func(error)) error {
	op := &taskData{
		event:       e,
		onValidated: onValidated,
	}
	select {
	case v.tasksQ <- op:
		return nil
	case <-v.quit:
		return errTerminated
	}
}

// verifySignature checks the signature against e.Creator.
func verifySignature(signedHash hash.Hash, sig inter.Signature, pubkey validatorpk.PubKey) bool {
	if pubkey.Type != validatorpk.Types.Secp256k1 {
		return false
	}
	return crypto.VerifySignature(pubkey.Raw, signedHash.Bytes(), sig.Bytes())
}

func (v *Checker) ValidateEventLocator(e inter.SignedEventLocator, authEpoch idx.Epoch, authErr error, checkPayload func() bool) error {
	pubkeys := v.reader.GetEpochPubKeysOf(authEpoch)
	if len(pubkeys) == 0 {
		return authErr
	}
	pubkey, ok := pubkeys[e.Locator.Creator]
	if !ok {
		return epochcheck.ErrAuth
	}
	if checkPayload != nil && !checkPayload() {
		return ErrWrongPayloadHash
	}
	if !verifySignature(e.Locator.HashToSign(), e.Sig, pubkey) {
		return ErrWrongEventSig
	}
	return nil
}

func (v *Checker) matchPubkey(creator idx.ValidatorID, epoch idx.Epoch, want []byte, authErr error) error {
	pubkeys := v.reader.GetEpochPubKeysOf(epoch)
	if len(pubkeys) == 0 {
		return authErr
	}
	pubkey, ok := pubkeys[creator]
	if !ok {
		return epochcheck.ErrAuth
	}
	if bytes.Compare(pubkey.Bytes(), want) != 0 {
		return ErrPubkeyChanged
	}
	return nil
}

// ValidateEvent runs heavy checks for event
func (v *Checker) ValidateEvent(e inter.EventPayloadI) error {
	pubkeys, epoch := v.reader.GetEpochPubKeys()
	if e.Epoch() != epoch {
		return epochcheck.ErrNotRelevant
	}
	// validatorID
	pubkey, ok := pubkeys[e.Creator()]
	if !ok {
		return epochcheck.ErrAuth
	}
	// event sig
	if !verifySignature(e.HashToSign(), e.Sig(), pubkey) {
		return ErrWrongEventSig
	}
	// pre-cache tx sig
	for _, tx := range e.Txs() {
		_, err := types.Sender(v.txSigner, tx)
		if err != nil {
			return ErrMalformedTxSig
		}
	}
	// Payload hash
	if e.PayloadHash() != inter.CalcPayloadHash(e) {
		return ErrWrongPayloadHash
	}

	if e.AnyBridgeVotes() {
		if !validateBridgeVotes(e.BridgeVotes(), pubkey) {
			return ErrInvalidBridgeVote
		}
	}

	return nil
}

func validateBridgeVotes(votes []inter.BridgeVote, pubkey validatorpk.PubKey) bool {
	for _, vote := range votes {
		var sign inter.BridgeSignature
		copy(sign[:], vote.Signature[:])
		if sign[64] < 27 || sign[64] > 28 {
			return false
		}
		sign[64] -= 27 // the V value is 27 or 28 as an Ethereum convention
		votePubkey, err := crypto.SigToPub(vote.Hash[:], sign[:])
		if err != nil {
			return false
		}
		if !bytes.Equal(crypto.FromECDSAPub(votePubkey), pubkey.Raw) {
			return false
		}
	}
	return true
}

func (v *Checker) loop() {
	defer v.wg.Done()
	for {
		select {
		case op := <-v.tasksQ:
			if op.event != nil {
				op.onValidated(v.ValidateEvent(op.event))
			}

		case <-v.quit:
			return
		}
	}
}
