package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"slices"
)

func hashItem(item string, nslots uint64) uint64 {
	// md5 -> 128 bit hash, endian-ess to preserve determinism, platform-agnostism,
	// XOR-ing to create a 64 bit output
	digest := md5.Sum([]byte(item))
	digestHigh := binary.BigEndian.Uint64(digest[8:16])
	digestLow := binary.BigEndian.Uint64(digest[:8])
	return (digestHigh ^ digestLow) % nslots
}

type ConsistentHasher struct {
	nodes  []string
	slots  []uint64
	nslots uint64
}

func NewConsistentHasher(nslots uint64) *ConsistentHasher {
	return &ConsistentHasher{
		nslots: nslots,
	}
}

func (ch *ConsistentHasher) AddNode(node string) error {
	if len(ch.nodes) >= int(ch.nslots) {
		return fmt.Errorf("Slots are full!")
	}

	nodeHash := hashItem(node, ch.nslots)
	slotPos, found := slices.BinarySearch(ch.slots, nodeHash)

	if found {
		return fmt.Errorf("Node collision!")
	}

	ch.slots = slices.Insert(ch.slots, slotPos, nodeHash)
	ch.nodes = slices.Insert(ch.nodes, slotPos, node)

	return nil
}

func (ch *ConsistentHasher) DeleteNode(node string) error {
	nodeHash := hashItem(node, ch.nslots)
	slotPos, found := slices.BinarySearch(ch.slots, nodeHash)

	if !found {
		return fmt.Errorf("Node doesn't exists!")
	}

	ch.slots = slices.Delete(ch.slots, slotPos, slotPos+1)
	ch.nodes = slices.Delete(ch.nodes, slotPos, slotPos+1)

	return nil
}

func (ch *ConsistentHasher) FindNodeFor(item string) string {
	itemHash := hashItem(item, ch.nslots)
	itemPos, _ := slices.BinarySearch(ch.slots, itemHash)

	if itemPos == len(ch.slots) {
		itemPos = 0
	}

	return ch.nodes[itemPos]

}
