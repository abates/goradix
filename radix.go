package goradix

import (
	"fmt"
)

var bitmasks = []uint64{
	0x0000000000000000, 0x8000000000000000, 0xc000000000000000, 0xe000000000000000, 0xf000000000000000,
	0xf800000000000000, 0xfc00000000000000, 0xfe00000000000000, 0xff00000000000000, 0xff80000000000000,
	0xffc0000000000000, 0xffe0000000000000, 0xfff0000000000000, 0xfff8000000000000, 0xfffc000000000000,
	0xfffe000000000000, 0xffff000000000000, 0xffff800000000000, 0xffffc00000000000, 0xffffe00000000000,
	0xfffff00000000000, 0xfffff80000000000, 0xfffffc0000000000, 0xfffffe0000000000, 0xffffff0000000000,
	0xffffff8000000000, 0xffffffc000000000, 0xffffffe000000000, 0xfffffff000000000, 0xfffffff800000000,
	0xfffffffc00000000, 0xfffffffe00000000, 0xffffffff00000000, 0xffffffff80000000, 0xffffffffc0000000,
	0xffffffffe0000000, 0xfffffffff0000000, 0xfffffffff8000000, 0xfffffffffc000000, 0xfffffffffe000000,
	0xffffffffff000000, 0xffffffffff800000, 0xffffffffffc00000, 0xffffffffffe00000, 0xfffffffffff00000,
	0xfffffffffff80000, 0xfffffffffffc0000, 0xfffffffffffe0000, 0xffffffffffff0000, 0xffffffffffff8000,
	0xffffffffffffc000, 0xffffffffffffe000, 0xfffffffffffff000, 0xfffffffffffff800, 0xfffffffffffffc00,
	0xfffffffffffffe00, 0xffffffffffffff00, 0xffffffffffffff80, 0xffffffffffffffc0, 0xffffffffffffffe0,
	0xfffffffffffffff0, 0xfffffffffffffff8, 0xfffffffffffffffc, 0xfffffffffffffffe, 0xffffffffffffffff,
}

type Node struct {
	prefix   uint64
	length   uint8
	children []*Node
	path     string
}

func (n *Node) IsLeaf() bool {
	return len(n.children) == 0
}

func (n *Node) String() string {
	return fmt.Sprintf("%2d   0x%016x   %064b", n.length, n.prefix, n.prefix)
}

func (n *Node) Insert(value *Node) {
	var match *Node
	maxLength := uint8(0)

	for _, child := range n.children {
		childLength := child.Match(value.prefix)
		if childLength > maxLength {
			match = child
			maxLength = childLength
		}
	}

	if match != nil {
		value.prefix = value.prefix << maxLength
		value.length -= maxLength

		if maxLength < match.length {
			newNode := &Node{
				prefix:   match.prefix << maxLength,
				length:   maxLength,
				children: match.children,
			}

			match.length = maxLength
			match.prefix = match.prefix & bitmasks[match.length]
			match.children = []*Node{newNode, value}
		} else {
			match.Insert(value)
		}
	} else {
		if n.length > 0 {
			n.children = append(n.children, &Node{
				prefix: n.prefix << (n.length - value.length),
				length: value.length,
			})

			n.length = n.length - value.length
			n.prefix = n.prefix & bitmasks[n.length]
		}

		n.children = append(n.children, value)
	}
}

func (n *Node) Match(value uint64) (length uint8) {
	b := uint64(0x8000000000000000)
	for l := uint8(0); l < n.length; l++ {
		if value&b != n.prefix&b {
			break
		}
		length += 1
		b = b >> 1
	}
	return
}

func (n *Node) Lookup(value uint64) *Node {
	var match *Node
	var maxLength uint8
	for _, child := range n.children {
		length := child.Match(value)
		if length > maxLength {
			match = child
			maxLength = length
		}
	}

	if match == nil {
		return n
	}

	return match.Lookup(value << maxLength)
}

type RadixTree struct {
	root *Node
}

func New() *RadixTree {
	r := new(RadixTree)
	r.root = new(Node)
	r.root.prefix = 0
	r.root.length = 0
	return r
}

func (r *RadixTree) Insert(value uint64) {
	node := &Node{
		prefix: value,
		length: 64,
	}
	r.root.Insert(node)
}

func (r *RadixTree) Lookup(value uint64) {
	r.root.Lookup(value)
}
