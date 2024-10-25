package tests

import (
	"testing"

	"github.com/Fantom-foundation/go-opera/tests/contracts/coinbase"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestCoinBase_CoinbaseYieldsZeroAddress(t *testing.T) {

	coinbaseAddress := common.Address{}

	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()

	contract, receipt, err := DeployContract(net, coinbase.DeployCoinbase)
	checkTxExecution(t, receipt, err)

	receipt, err = net.Apply(contract.LogCoinBaseAddress)
	checkTxExecution(t, receipt, err)

	if len(receipt.Logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(receipt.Logs))
	}
	if want, got := coinbaseAddress, common.BytesToAddress(receipt.Logs[0].Data); want != got {
		t.Errorf("Expected coinbase address %v, got %v", want, got)
	}
}

func TestCoinBase_CoinbaseAccessIsAlwaysWarm(t *testing.T) {

	someAccountAddress := common.Address{1}
	coinbaseAddress := common.Address{}
	someBalance := int64(1234)

	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()

	contract, receipt, err := DeployContract(net, coinbase.DeployCoinbase)
	checkTxExecution(t, receipt, err)

	// Create some account by transferring balance
	receipt, err = net.EndowAccount(someAccountAddress, someBalance)
	checkTxExecution(t, receipt, err)
	// Create the coinbase account (if not, creation fees will be charged)
	receipt, err = net.EndowAccount(coinbaseAddress, someBalance)
	checkTxExecution(t, receipt, err)

	// touch account (COLD)
	receipt, err = net.Apply(func(ops *bind.TransactOpts) (*types.Transaction, error) {
		ops.GasPrice = nil                  // setting this field forces legacy-tx use
		ops.AccessList = types.AccessList{} // leave access list empty to force same kind of transaction
		return contract.TouchAddress(ops, someAccountAddress)
	})
	checkTxExecution(t, receipt, err)
	coldCost := receipt.GasUsed

	// touch account (WARM)
	receipt, err = net.Apply(func(ops *bind.TransactOpts) (*types.Transaction, error) {
		ops.GasPrice = nil // setting this field forces legacy-tx use
		ops.AccessList = types.AccessList{
			// warm access of account
			{Address: someAccountAddress, StorageKeys: []common.Hash{}},
		}
		return contract.TouchAddress(ops, someAccountAddress)
	})
	checkTxExecution(t, receipt, err)
	warmCost := receipt.GasUsed

	// Touch coinbase
	receipt, err = net.Apply(func(ops *bind.TransactOpts) (*types.Transaction, error) {
		ops.GasPrice = nil                  // setting this field forces legacy-tx use
		ops.AccessList = types.AccessList{} // leave access list empty to force same kind of transaction
		return contract.TouchCoinbase(ops)
	})

	checkTxExecution(t, receipt, err)
	coinBaseCost := receipt.GasUsed

	// Difference between cold and warm access point out the amount of gas burned (plus cold access)
	transactionCostOverhead := coldCost - warmCost
	coinbaseDiff := (coldCost - coinBaseCost) - transactionCostOverhead
	// Remove the extra cost of a cold access
	differenceBetweenContracts := coinbaseDiff - 2500

	// The instructions executed by the contract differ slightly between the two calls.
	// Static costs and small memory copies account for a small difference in gas consumed.
	gasDelta := 80
	if differenceBetweenContracts > uint64(gasDelta) {
		t.Errorf("Expected difference between gas consumed by contracts to be less than %d, got %d", gasDelta, differenceBetweenContracts)
	}
}

////////////////////////////////////////////////////////////////////////////////
// helpers

func checkTxExecution(t *testing.T, receipt *types.Receipt, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("failed to execute transaction; %v", err)
	}
	if want, got := types.ReceiptStatusSuccessful, receipt.Status; want != got {
		t.Errorf("Expected status %v, got %v", want, got)
	}
}
