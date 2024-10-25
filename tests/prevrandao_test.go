package tests

import (
	"context"
	"github.com/Fantom-foundation/go-opera/tests/contracts/prevrandao"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"testing"
)

func TestPrevRandao(t *testing.T) {
	net, err := StartIntegrationTestNet(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to start the fake network: %v", err)
	}
	defer net.Stop()

	// Deploy the contract.
	contract, _, err := DeployContract(net, prevrandao.DeployPrevrandao)
	if err != nil {
		t.Fatalf("failed to deploy contract; %v", err)
	}

	// Collect the current PrevRandao fee from the head state.
	receipt, err := net.Apply(contract.LogCurrentPrevRandao)
	if err != nil {
		t.Fatalf("failed to log current base fee; %v", err)
	}

	if len(receipt.Logs) != 1 {
		t.Fatalf("unexpected number of logs; expected 1, got %d", len(receipt.Logs))
	}

	entry, err := contract.ParseCurrentPrevRandao(*receipt.Logs[0])
	if err != nil {
		t.Fatalf("failed to parse log; %v", err)
	}
	fromLog := entry.Prevrandao

	// Collect the prevrandao from the block header.
	client, err := net.GetClient()
	if err != nil {
		t.Fatalf("failed to get client; %v", err)
	}
	defer client.Close()

	block, err := client.BlockByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		t.Fatalf("failed to get block header; %v", err)
	}
	fromBlock := block.MixDigest().Big() // MixDigest == MixHash == PrevRandao

	// Collect the base fee from the archive.
	fromArchive, err := contract.GetPrevRandao(&bind.CallOpts{BlockNumber: receipt.BlockNumber})
	if err != nil {
		t.Fatalf("failed to get base fee from archive; %v", err)
	}

	if fromLog.Sign() < 1 {
		t.Fatalf("invalid base fee from log; %v", fromLog)
	}

	if fromLog.Cmp(fromBlock) != 0 {
		t.Fatalf("base fee mismatch; from log %v, from block %v", fromLog, fromBlock)
	}
	if fromLog.Cmp(fromArchive) != 0 {
		t.Fatalf("base fee mismatch; from log %v, from archive %v", fromLog, fromArchive)
	}
}
