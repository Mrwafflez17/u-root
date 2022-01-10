// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/ls"
)

type ttflags struct {
	all       bool
	human     bool
	directory bool
	long      bool
	quoted    bool
	recurse   bool
	classify  bool
}

// Test listName func
func TestListName(t *testing.T) {
	// Create some directorys.
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "testDir")
	os.Mkdir(testDir, 0o700)
	os.Mkdir(filepath.Join(testDir, "d1"), 0o740)
	// Create some files.
	files := []string{"f1", "f2", "f3\nline 2", ".f4", "d1/f4"}
	for _, file := range files {
		os.Create(filepath.Join(testDir, file))
	}
	defer os.RemoveAll(tmpDir)

	// Creating test table
	for _, tt := range []struct {
		name   string
		input  string
		want   string
		flag   ttflags
		prefix bool
	}{
		{
			name:  "ls without arguments",
			input: testDir,
			want:  fmt.Sprintf("%s\n%s\n%s\n%s\n", "d1", "f1", "f2", "f3?line 2"),
			flag:  ttflags{},
		},
		{
			name:  "ls osfi.IsDir() path, quoted = true",
			input: testDir,
			want:  fmt.Sprintf("\"%s\"\n\"%s\"\n\"%s\"\n\"%s\"\n", "d1", "f1", "f2", "f3\\nline 2"),
			flag: ttflags{
				quoted: true,
			},
		},
		{
			name:  "ls osfi.IsDir() path, quoted = true, prefix = true ",
			input: testDir,
			want:  fmt.Sprintf("\"%s\":\n\"%s\"\n\"%s\"\n\"%s\"\n\"%s\"\n", testDir, "d1", "f1", "f2", "f3\\nline 2"),
			flag: ttflags{
				quoted: true,
			},
			prefix: true,
		},
		{
			name:   "ls osfi.IsDir() path, quoted = false, prefix = true ",
			input:  testDir,
			want:   fmt.Sprintf("%s:\n%s\n%s\n%s\n%s\n", testDir, "d1", "f1", "f2", "f3?line 2"),
			prefix: true,
		},
		{
			name:  "ls recurse = true",
			input: testDir,
			want: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n", testDir, filepath.Join(testDir, ".f4"), filepath.Join(testDir, "d1"),
				filepath.Join(testDir, "d1/f4"), filepath.Join(testDir, "f1"), filepath.Join(testDir, "f2"), filepath.Join(testDir, "f3?line 2")),
			flag: ttflags{
				all:     true,
				recurse: true,
			},
		},
		{
			name:  "ls directory = true",
			input: testDir,
			want:  fmt.Sprintf("%s\n", "testDir"),
			flag: ttflags{
				directory: true,
			},
		},
		{
			name:  "ls classify = true",
			input: testDir,
			want:  fmt.Sprintf("%s\n%s\n%s\n%s\n", "d1/", "f1", "f2", "f3?line 2"),
			flag: ttflags{
				classify: true,
			},
		},
		{
			name:  "file does not exist",
			input: "dir",
			flag:  ttflags{},
		},
	} {
		// Setting the flags
		*all = tt.flag.all
		*human = tt.flag.human
		*directory = tt.flag.directory
		*long = tt.flag.long
		*quoted = tt.flag.quoted
		*recurse = tt.flag.recurse
		*classify = tt.flag.classify

		// Running the tests
		t.Run(tt.name, func(t *testing.T) {
			// Write output in buffer.
			var buf bytes.Buffer
			var s ls.Stringer = ls.NameStringer{}
			if *quoted {
				s = ls.QuotedStringer{}
			}
			if *long {
				s = ls.LongStringer{Human: *human, Name: s}
			}
			if err := listName(s, tt.input, &buf, tt.prefix); err != nil {
				if buf.String() != tt.want {
					t.Errorf("listName() = '%v', want: '%v'", buf.String(), tt.want)
				}
			}
		})
	}
}

// Test list func
func TestList(t *testing.T) {
	// Create some directorys.
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)

	// Creating test table
	for _, tt := range []struct {
		name   string
		input  []string
		want   error
		flag   ttflags
		prefix bool
	}{
		{
			name:  "input empty, quoted = true, long = true",
			input: []string{},
			want:  nil,
			flag: ttflags{
				quoted: true,
				long:   true,
			},
		},
		{
			name:  "input empty, quoted = true, long = true",
			input: []string{"dir"},
			want:  fmt.Errorf("error while listing %v: 'lstat %v: no such file or directory'", "dir", "dir"),
		},
	} {

		// Setting the flags
		*long = tt.flag.long
		*quoted = tt.flag.quoted

		// Running the tests
		t.Run(tt.name, func(t *testing.T) {
			if got := list(tt.input); got != nil {
				if got.Error() != tt.want.Error() {
					t.Errorf("list() = '%v', want: '%v'", got, tt.want)
				}
			}
		})
	}
}

// Test indicator func
func TestIndicator(t *testing.T) {
	// Creating test table
	for _, test := range []struct {
		lsInfo ls.FileInfo
		symbol string
	}{
		{
			ls.FileInfo{
				Mode: os.ModeDir,
			},
			"/",
		},
		{
			ls.FileInfo{
				Mode: os.ModeNamedPipe,
			},
			"|",
		},
		{
			ls.FileInfo{
				Mode: os.ModeSymlink,
			},
			"@",
		},
		{
			ls.FileInfo{
				Mode: os.ModeSocket,
			},
			"=",
		},
		{
			ls.FileInfo{
				Mode: 0b110110100,
			},
			"",
		},
		{
			ls.FileInfo{
				Mode: 0b111111101,
			},
			"*",
		},
	} {
		// Run tests
		got := indicator(test.lsInfo)
		if got != test.symbol {
			t.Errorf("for mode '%b' expected '%q', got '%q'", test.lsInfo.Mode, test.symbol, got)
		}
	}
}
