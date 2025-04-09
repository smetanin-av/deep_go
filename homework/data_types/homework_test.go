package main

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func ToLittleEndian[T ~uint16 | ~uint32 | ~uint64 | ~uint](val T) T {
	var res T
	for range unsafe.Sizeof(val) {
		res = (res << 8) | (val & 0xFF)
		val >>= 8
	}
	return res
}

func TestСonversionUint16(t *testing.T) {
	for _, tc := range []struct {
		val, exp uint16
	}{
		{val: 0x00_00, exp: 0x00_00},
		{val: 0xFF_FF, exp: 0xFF_FF},
		{val: 0x00_FF, exp: 0xFF_00},
		{val: 0x01_02, exp: 0x02_01},
		{val: 0x0F_F0, exp: 0xF0_0F},
	} {
		t.Run(fmt.Sprintf("%04X should be %04X", tc.val, tc.exp), func(t *testing.T) {
			res := ToLittleEndian(tc.val)
			assert.Equal(t, tc.exp, res)
		})
	}
}

func TestСonversionUint32(t *testing.T) {
	for _, tc := range []struct {
		val, exp uint32
	}{
		{val: 0x00_00_00_00, exp: 0x00_00_00_00},
		{val: 0xFF_FF_FF_FF, exp: 0xFF_FF_FF_FF},
		{val: 0x00_FF_00_FF, exp: 0xFF_00_FF_00},
		{val: 0x00_00_FF_FF, exp: 0xFF_FF_00_00},
		{val: 0x01_02_03_04, exp: 0x04_03_02_01},
		{val: 0x0F_F0_0F_F0, exp: 0xF0_0F_F0_0F},
	} {
		t.Run(fmt.Sprintf("%08X should be %08X", tc.val, tc.exp), func(t *testing.T) {
			res := ToLittleEndian(tc.val)
			assert.Equal(t, tc.exp, res)
		})
	}
}

func TestСonversionUint64(t *testing.T) {
	for _, tc := range []struct {
		val, exp uint64
	}{
		{val: 0x00_00_00_00_00_00_00_00, exp: 0x00_00_00_00_00_00_00_00},
		{val: 0xFF_FF_FF_FF_FF_FF_FF_FF, exp: 0xFF_FF_FF_FF_FF_FF_FF_FF},
		{val: 0x00_FF_00_FF_00_FF_00_FF, exp: 0xFF_00_FF_00_FF_00_FF_00},
		{val: 0x00_00_FF_FF_00_00_FF_FF, exp: 0xFF_FF_00_00_FF_FF_00_00},
		{val: 0x01_02_03_04_05_06_07_08, exp: 0x08_07_06_05_04_03_02_01},
		{val: 0x0F_F0_0F_F0_0F_F0_0F_F0, exp: 0xF0_0F_F0_0F_F0_0F_F0_0F},
	} {
		t.Run(fmt.Sprintf("%016X should be %016X", tc.val, tc.exp), func(t *testing.T) {
			res := ToLittleEndian(tc.val)
			assert.Equal(t, tc.exp, res)
		})
	}
}
