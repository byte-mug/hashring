/*
MIT License

Copyright (c) 2017 Simon Schmidt
Copyright (c) 2016 Sung-jin Hong

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/


package hashring

import "math"
import "sort"
import "crypto/md5"

type uints []uint32
func (h uints) Len() int           { return len(h) }
func (h uints) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h uints) Less(i, j int) bool { return h[i] < h[j] }

type Binary struct {
	data []byte
}
func NewBinary(b []byte) Binary { return Binary{b} }
func (b Binary) Md5() [md5.Size]byte { return md5.Sum(b.data) }
func (b Binary) String() string { return string(b.data) }
func (b Binary) ToBytes(prefix []byte) []byte { return append(prefix,b.data...) }

type RingNode struct {
	Key Binary
	Value interface{}
	Weigth int
}

type HashRing struct {
	Table []RingNode
	index map[uint32]int
	keys  uints
}
// This code is partially derived from https://github.com/serialx/hashring
func (h *HashRing) GenerateCircle() {
	var m5h [md5.Size]byte
	h.index = make(map[uint32]int)
	h.keys = h.keys[:0]
	totalWeight := 0
	for _,v := range h.Table {
		weigth := v.Weigth
		if weigth<1 { weigth = 1 }
		totalWeight += weigth
	}
	
	for i,v := range h.Table {
		weigth := v.Weigth
		if weigth<1 { weigth = 1 }
		
		factor := int(math.Floor(float64(40*len(h.Table)*weigth) / float64(totalWeight)))
		m5h = v.Key.Md5()
		k := crcMd5(0,m5h)
		
		for j:=0 ; j<factor ; j++ {
			k2 := uint32(crcInt(k,j))
			
			h.index[k2] = i
			h.keys = append(h.keys,k2)
		}
	}
	
	sort.Sort(h.keys)
}

func encodeLE(b []byte,d int) {
	for i := range b {
		b[i] = byte(d)
		d >>= 8
	}
}

func (h *HashRing) GetNodePosition(b Binary) int {
	return h.GetNodePositionIterative(b,0)
}

func (h *HashRing) GetNodePositionIterative(b Binary, iteration int) int {
	m5h := b.Md5()
	mcrc := crcMd5(0,m5h)
	hkey := uint32(crcInt(mcrc,iteration))
	keys := h.keys
	if len(keys)==0 { return -1 }
	
	position := sort.Search(len(keys), func(i int) bool { return keys[i] > hkey } )
	
	if position >= len(keys) { position = 0 }
	return h.index[keys[position]]
}

func (h *HashRing) GetAllNodePositions(b Binary, out []int, begin int) bool {
	m5h := b.Md5()
	mcrc := crcMd5(0,m5h)
	hkey := uint32(0)
	keys := h.keys
	if len(keys)==0 { return false }
	
	lk := len(keys)
	find := func(i int) bool { return keys[i] > hkey }
	
	for i := range out {
		hkey = uint32(crcInt(mcrc,begin))
		begin++
		
		position := sort.Search(lk, find)
		
		if position >= len(keys) { position = 0 }
		out[i] = h.index[keys[position]]
	}
	return true
}

