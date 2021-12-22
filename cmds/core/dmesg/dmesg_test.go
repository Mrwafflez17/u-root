// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestDmesg(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	for _, tt := range []struct {
		name      string
		clear     bool
		readClear bool
		want      error
	}{
		{
			name:      "both flags set",
			clear:     true,
			readClear: true,
			want:      fmt.Errorf("cannot specify both -clear and -read-clear"),
		},
		{
			name:      "clear log",
			clear:     true,
			readClear: false,
			want:      fmt.Errorf(""),
		},
		{
			name:      "clear log after printing",
			clear:     false,
			readClear: true,
			want:      fmt.Errorf(""),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			*clear = tt.clear
			*readClear = tt.readClear
			buf := &bytes.Buffer{}
			if got := dmesg(buf); got != nil {
				if got.Error() != tt.want.Error() {
					t.Errorf("dmesg() = '%v', want: '%v'", got, tt.want)
				}
			} else {
				if buf.String() != "" && (*clear || *readClear) {
					t.Errorf("System log should be cleared")
				}
			}
		})
	}
}
