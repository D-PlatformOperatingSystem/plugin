// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mm

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/holiman/uint256"
)

// Memory       ， EVM            
type Memory struct {
	// Store         
	Store []byte
	// LastGasCost          Gas
	LastGasCost uint64
}

// NewMemory         
func NewMemory() *Memory {
	return &Memory{}
}

// Set        ， value => offset:offset + size
func (m *Memory) Set(offset, size uint64, value []byte) (err error) {
	if size > 0 {
		//    +            
		if offset+size > uint64(len(m.Store)) {
			err = fmt.Errorf("INVALID memory access, memory size:%v, offset:%v, size:%v", len(m.Store), offset, size)
			log15.Crit(err.Error())
			//panic("invalid memory: store empty")
			return err
		}
		copy(m.Store[offset:offset+size], value)
	}
	return nil
}

// Set32  offset    32       ，       32   ，     
func (m *Memory) Set32(offset uint64, val *uint256.Int) (err error) {

	//          
	if offset+32 > uint64(len(m.Store)) {
		err = fmt.Errorf("INVALID memory access, memory size:%v, offset:%v, size:%v", len(m.Store), offset, 32)
		log15.Crit(err.Error())
		//panic("invalid memory: store empty")
		return err
	}
	//      
	copy(m.Store[offset:offset+32], []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	// Fill in relevant bits
	val.WriteToSlice(m.Store[offset:])

	return nil
}

// Resize          
func (m *Memory) Resize(size uint64) {
	if uint64(m.Len()) < size {
		m.Store = append(m.Store, make([]byte, size-uint64(m.Len()))...)
	}
}

// Get                     ，           
func (m *Memory) Get(offset, size int64) (cpy []byte) {
	if size == 0 {
		return nil
	}

	if len(m.Store) > int(offset) {
		cpy = make([]byte, size)
		copy(cpy, m.Store[offset:offset+size])

		return
	}

	return
}

// GetPtr  Get  ，            
func (m *Memory) GetPtr(offset, size int64) []byte {
	if size == 0 {
		return nil
	}

	if len(m.Store) > int(offset) {
		return m.Store[offset : offset+size]
	}

	return nil
}

// Len              （     ）
func (m *Memory) Len() int {
	return len(m.Store)
}

// Data             
func (m *Memory) Data() []byte {
	return m.Store
}

// Print         （   ）
func (m *Memory) Print() {
	fmt.Printf("### mem %d bytes ###\n", len(m.Store))
	if len(m.Store) > 0 {
		addr := 0
		for i := 0; i+32 <= len(m.Store); i += 32 {
			fmt.Printf("%03d: % x\n", addr, m.Store[i:i+32])
			addr++
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("####################")
}

//                ，        
func calcMemSize64(off, l *uint256.Int) (uint64, bool) {
	if !l.IsUint64() {
		return 0, true
	}
	return calcMemSize64WithUint(off, l.Uint64())
}

// calcMemSize64WithUint calculates the required memory size, and returns
// the size and whether the result overflowed uint64
// Identical to calcMemSize64, but length is a uint64
func calcMemSize64WithUint(off *uint256.Int, length64 uint64) (uint64, bool) {
	// if length is zero, memsize is always zero, regardless of offset
	if length64 == 0 {
		return 0, false
	}
	// Check that offset doesn't overflow
	offset64, overflow := off.Uint64WithOverflow()
	if overflow {
		return 0, true
	}
	val := offset64 + length64
	// if value < either of it's parts, then it overflowed
	return val, val < offset64
}
