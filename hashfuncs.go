/*
MIT License

Copyright (c) 2017 Simon Schmidt

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

import "hash"
import "hash/fnv"
import "hash/crc64"
import "sync"

var pool_Fnv = sync.Pool{ New : func() interface{} { return fnv.New64a() } }

func Hash(b []byte) uint64 {
	h := pool_Fnv.Get().(hash.Hash64)
	defer pool_Fnv.Put(h)
	h.Reset()
	h.Write(b)
	return h.Sum64()
}

var crctable = crc64.MakeTable(crc64.ECMA)

func Derive(base uint64,b []byte) uint64 {
	return crc64.Update(base,crctable,b)
}

