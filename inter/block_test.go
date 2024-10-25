package inter

import (
	"github.com/Fantom-foundation/lachesis-base/hash"
	"math/rand"
	"reflect"
	"testing"
)

func TestBlock_ComputePrevRandao_ComputationIsDeterministic(t *testing.T) {
	events := []hash.Event{
		{byte(rand.Int())},
		{byte(rand.Int())},
		{byte(rand.Int())},
		{byte(rand.Int())},
	}
	blk := Block{Events: events}
	randao1 := blk.ComputePrevRandao()
	rand.Shuffle(len(blk.Events), func(i, j int) {
		blk.Events[i], blk.Events[j] = blk.Events[j], blk.Events[i]
	})
	randao2 := blk.ComputePrevRandao()
	if !reflect.DeepEqual(randao1, randao2) {
		t.Error("computation is not deterministic")
	}
}

func TestBlock_ComputePrevRandao_First8BytesAreIgnored(t *testing.T) {
	b := make([]byte, 32)
	for i := range 8 {
		b[i] = byte(i)
	}
	blk := Block{Events: []hash.Event{
		hash.BytesToEvent(b),
	}}
	randao := blk.ComputePrevRandao()
	for i := range 8 {
		if randao[i] != 0 {
			t.Errorf("byte at position %d is not zero", i)
		}
	}
}
