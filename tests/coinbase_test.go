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

	// touch account, with address in access list (WARM)
	receipt, err = net.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		opts.GasPrice = nil // setting this field forces legacy-tx use
		opts.AccessList = types.AccessList{
			// warm access of account
			{Address: someAccountAddress, StorageKeys: []common.Hash{}},
		}
		return contract.TouchAddress(opts, someAccountAddress)
	})
	checkTxExecution(t, receipt, err)
	warmCost := receipt.GasUsed

	// Get coinbase address
	receipt, err = net.Apply(contract.LogCoinBaseAddress)
	checkTxExecution(t, receipt, err)
	if len(receipt.Logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(receipt.Logs))
	}
	coinbaseAddress := common.BytesToAddress(receipt.Logs[0].Data)

	// Touch coinbase without access list (COLD)
	receipt, err = net.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return contract.TouchAddress(opts, coinbaseAddress)
	})
	checkTxExecution(t, receipt, err)
	coinBaseCost := receipt.GasUsed

	// CoinBase address access has the same cost as warm access
	if warmCost == coinBaseCost {
		t.Errorf("Expected coinbase access to be cheaper than warm access")
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
