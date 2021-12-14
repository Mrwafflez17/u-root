// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmos

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/memio"
)

func getMock(err string, inBuf, outBuf io.ReadWriter) *CMOSChip {
	return &CMOSChip{
		In: func(addr uint16, data memio.UintN) error {
			if err != "" {
				return fmt.Errorf(err)
			}
			switch addr {
			case cmosRegPort:
				return nil
			case cmosDataPort:
				io.Copy(inBuf, outBuf)
				return nil
			default:
				return fmt.Errorf("invalid address")
			}
		},
		Out: func(addr uint16, data memio.UintN) error {
			if err != "" {
				return fmt.Errorf(err)
			}
			switch addr {
			case cmosRegPort:
				return nil
			case cmosDataPort:
				outBuf.Write([]byte(data.String()))
				return nil
			default:
				return fmt.Errorf("invalid address")
			}
		},
	}
}

func TestCMOS(t *testing.T) {
	for _, tt := range []struct {
		name                string
		addr                memio.Uint8
		writeData, readData memio.UintN
		err                 string
	}{
		{
			name:      "uint8",
			addr:      0x10,
			writeData: &[]memio.Uint8{0x12}[0],
		},
		{
			name:      "uint16",
			addr:      0x20,
			writeData: &[]memio.Uint16{0x1234}[0],
		},
		{
			name:      "uint32",
			addr:      0x30,
			writeData: &[]memio.Uint32{0x12345678}[0],
		},
		{
			name:      "uint64",
			addr:      0x40,
			writeData: &[]memio.Uint64{0x1234567890abcdef}[0],
			err:       "data must not be of type Uint64",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var in, out bytes.Buffer
			c := getMock(tt.err, &in, &out)
			// Set internal function to dummy but save old state for reset later
			if err := c.Write(tt.addr, tt.writeData); err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Error(err)
				}
			}
			err := c.Read(tt.addr, tt.readData)
			if err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Error(err)
				}
			}
			// We can only progress if error is nil.
			if err == nil {
				got := in.String()
				want := tt.writeData.String()
				if got != want {
					t.Errorf("Got %s, want %s", got, want)
				}
			}
		})
	}
}

// This is just for coverage percentage. This test does nothing of any other value.
func TestGetCMOS(t *testing.T) {
	_ = GetCMOS()
}
