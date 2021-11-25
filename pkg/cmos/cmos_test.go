// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmos

import (
	"errors"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/memio"
)

var tests = []struct {
	name                string
	addr                memio.Uint8
	writeData, readData memio.UintN
	err                 string
}{
	{
		name:      "uint8",
		addr:      0x10,
		writeData: &[]memio.Uint8{0x12}[0],
		readData:  new(memio.Uint8),
	},
	{
		name:      "uint16",
		addr:      0x20,
		writeData: &[]memio.Uint16{0x1234}[0],
		readData:  new(memio.Uint16),
	},
	{
		name:      "uint32",
		addr:      0x30,
		writeData: &[]memio.Uint32{0x12345678}[0],
		readData:  new(memio.Uint32),
	},
	{
		name:      "uint64",
		addr:      0x40,
		writeData: &[]memio.Uint64{0x1234567890abcdef}[0],
		readData:  new(memio.Uint64),
		err:       "data must not be of type Uint64",
	},
}

func TestCMOS(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set internal function to dummy but save old state for reset later
			oMemioOut := memioOut
			memioOut = func(addr uint16, data memio.UintN) (err error) {
				// Comment in memio/port_linux.go says Uint64 is not allowed, but checks for data.(*Uint8)
				// ToDo: Check what's right
				if data.Size() != int64(1) {
					return errors.New(tt.err)
				}
				return nil
			}
			defer func() { memioOut = oMemioOut }()

			// Set internal function to dummy but save old state for reset later
			oMemioIn := memioIn
			memioIn = func(addr uint16, data memio.UintN) (err error) { tt.readData = tt.writeData; return nil }
			defer func() { memioIn = oMemioIn }()

			if err := Write(tt.addr, tt.writeData); err != nil {
				if err.Error() != tt.err {
					t.Error(err)
				}
			}
			if err := Read(tt.addr, tt.readData); err != nil {
				if err.Error() != tt.err {
					t.Error(err)
				}
			}

			got := tt.writeData
			want := tt.readData
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Got %v, want %v", got, want)
			}
		})
	}
}
