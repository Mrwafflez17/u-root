// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64 || 386
// +build amd64 386

package cmos

import (
	"github.com/u-root/u-root/pkg/memio"
)

const (
	cmosRegPort  = 0x70
	cmosDataPort = 0x71
)

type CMOSChip struct {
	In  func(uint16, memio.UintN) error
	Out func(uint16, memio.UintN) error
}

// Read reads a register reg from CMOS into data.
func (c *CMOSChip) Read(reg memio.Uint8, data memio.UintN) error {
	if err := c.In(cmosRegPort, &reg); err != nil {
		return err
	}
	return c.In(cmosDataPort, data)
}

// Write writes value data into CMOS register reg.
func (c *CMOSChip) Write(reg memio.Uint8, data memio.UintN) error {
	if err := c.Out(cmosRegPort, &reg); err != nil {
		return err
	}
	return c.Out(cmosDataPort, data)
}

// GetCMOS() returns the struct to call Read and Write functions for CMOS
// associated with the correct functions of memio.In and memio.Out
func GetCMOS() *CMOSChip {
	return &CMOSChip{
		In:  memio.In,
		Out: memio.Out,
	}
}
