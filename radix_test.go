package goradix

import (
	"testing"
)

func TestMatch(t *testing.T) {
	n := &Node{
		prefix: 0xe900000000000000,
		length: 7,
	}

	v := uint64(0xef00000000000000)
	length := n.Match(v)

	if length != 5 {
		t.Errorf("Expected length 5 got %d", length)
	}
}

func checkValue(t *testing.T, n *Node, expectedValue uint64) {
	if n.prefix != expectedValue {
		t.Errorf("Expected 0x%016x, but got 0x%016x", expectedValue, n.prefix)
	}
}

func TestInsert(t *testing.T) {
	r := New()
	r.Insert(0x4a00000000000000)

	if len(r.root.children) < 1 {
		t.Errorf("Expected the new value to be inserted into the root node")
	}

	checkValue(t, r.root.children[0], 0x4a00000000000000)

	r.Insert(0x5d00000000000000)

	checkValue(t, r.root.children[0], 0x4000000000000000)

	if len(r.root.children[0].children) != 2 {
		t.Errorf("Expected the subtree to have two nodes but it has %d", len(r.root.children[0].children))
	}

	checkValue(t, r.root.children[0].children[0], 0x5000000000000000)
	checkValue(t, r.root.children[0].children[1], 0xe800000000000000)

	r.Insert(0x5900000000000000)

	checkValue(t, r.root.children[0].children[1], 0xc000000000000000)
	checkValue(t, r.root.children[0].children[1].children[0], 0xa000000000000000)
	checkValue(t, r.root.children[0].children[1].children[1], 0x2000000000000000)

	r.Insert(0x6900000000000000)
	checkValue(t, r.root.children[0], 0x4000000000000000)
	checkValue(t, r.root.children[0].children[0], 0x0000000000000000)

	checkValue(t, r.root.children[0].children[0].children[0], 0x5000000000000000)
	checkValue(t, r.root.children[0].children[0].children[1], 0xc000000000000000)
	checkValue(t, r.root.children[0].children[0].children[1].children[0], 0xa000000000000000)
	checkValue(t, r.root.children[0].children[0].children[1].children[1], 0x2000000000000000)

	checkValue(t, r.root.children[0].children[1], 0xa400000000000000)
}
