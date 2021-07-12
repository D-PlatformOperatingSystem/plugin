// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mm

import (
	"fmt"
	"sync"

	"github.com/holiman/uint256"
)

var stackPool = sync.Pool{
	New: func() interface{} {
		return &Stack{data: make([]uint256.Int, 0, 16)}
	},
}

// Stack      ，
type Stack struct {
	data []uint256.Int
}

// NewStack
func NewStack() *Stack {
	return stackPool.Get().(*Stack)
}

// Returnstack     stack  stackpool
func Returnstack(s *Stack) {
	s.data = s.data[:0]
	stackPool.Put(s)
}

// Data
func (st *Stack) Data() []uint256.Int {
	return st.data
}

// Push
func (st *Stack) Push(d *uint256.Int) {
	// NOTE push limit (1024) is checked in baseCheck
	st.data = append(st.data, *d)
}

// PushN
func (st *Stack) PushN(ds ...uint256.Int) {
	// FIXME: Is there a way to pass args by pointers.
	st.data = append(st.data, ds...)
}

// Pop
func (st *Stack) Pop() (ret uint256.Int) {
	ret = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-1]
	return
}

// Len
func (st *Stack) Len() int {
	return len(st.data)
}

// Swap
func (st *Stack) Swap(n int) {
	st.data[st.Len()-n], st.data[st.Len()-1] = st.data[st.Len()-1], st.data[st.Len()-n]
}

// Dup
func (st *Stack) Dup(n int) {
	st.Push(&st.data[st.Len()-n])
}

// Peek
func (st *Stack) Peek() *uint256.Int {
	return &st.data[st.Len()-1]
}

// Back    n
func (st *Stack) Back(n int) *uint256.Int {
	return &st.data[st.Len()-n-1]
}

// Require
func (st *Stack) Require(n int) error {
	if st.Len() < n {
		return fmt.Errorf("stack underflow (%d <=> %d)", len(st.data), n)
	}
	return nil
}

// Print     （   ）
func (st *Stack) Print() {
	fmt.Println("### stack ###")
	if len(st.data) > 0 {
		for i, val := range st.data {
			fmt.Printf("%-3d  %v\n", i, val)
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("#############")
}

var rStackPool = sync.Pool{
	New: func() interface{} {
		return &ReturnStack{data: make([]uint32, 0, 10)}
	},
}

// ReturnStack
type ReturnStack struct {
	data []uint32
}

// NewReturnStack        ，
func NewReturnStack() *ReturnStack {
	return rStackPool.Get().(*ReturnStack)
}

// ReturnRStack  returnStack  rStackPool
func ReturnRStack(rs *ReturnStack) {
	rs.data = rs.data[:0]
	rStackPool.Put(rs)
}

// Push
func (st *ReturnStack) Push(d uint32) {
	st.data = append(st.data, d)
}

// Pop  A uint32 is sufficient as for code below 4.2G
func (st *ReturnStack) Pop() (ret uint32) {
	ret = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-1]
	return
}

// Len ReturnStack
func (st *ReturnStack) Len() int {
	return len(st.data)
}

// Data
func (st *ReturnStack) Data() []uint32 {
	return st.data
}
