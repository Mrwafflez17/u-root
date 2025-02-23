// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/u-root/u-root/pkg/bb"
	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/ulog"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

// Commands to skip building in bb mode.
var skip = map[string]struct{}{
	"bb": {},
}

// BBBuilder is an implementation of Builder that compiles many Go commands
// into one busybox-style binary.
//
// BBBuilder will also include symlinks for each command to the busybox binary.
//
// BBBuilder does all this by rewriting the source files of the packages given
// to create one busybox-like binary containing all commands.
//
// The compiled binary uses argv[0] to decide which Go command to run.
//
// See bb/README.md for a detailed explanation of the implementation of busybox
// mode.
type BBBuilder struct {
	// ShellBang means generate #! files instead of symlinks.
	// ShellBang are more portable and just as efficient.
	ShellBang bool
}

// DefaultBinaryDir implements Builder.DefaultBinaryDir.
//
// The default initramfs binary dir is bbin for busybox binaries.
func (BBBuilder) DefaultBinaryDir() string {
	return "bbin"
}

// Build is an implementation of Builder.Build for a busybox-like initramfs.
func (b BBBuilder) Build(l ulog.Logger, af *initramfs.Files, opts Opts) error {
	// Build the busybox binary.
	bbPath := filepath.Join(opts.TempDir, "bb")
	noStrip := false
	if opts.BuildOpts != nil {
		noStrip = opts.BuildOpts.NoStrip
	}
	if err := bb.BuildBusybox(opts.Env, opts.Packages, noStrip, bbPath); err != nil {
		return err
	}

	if len(opts.BinaryDir) == 0 {
		return fmt.Errorf("must specify binary directory")
	}

	if err := af.AddFile(bbPath, path.Join(opts.BinaryDir, "bb")); err != nil {
		return err
	}

	// Add symlinks for included commands to initramfs.
	for _, pkg := range opts.Packages {
		if _, ok := skip[path.Base(pkg)]; ok {
			continue
		}

		// Add a symlink /bbin/{cmd} -> /bbin/bb to our initramfs.
		// Or add a #! file if b.ShellBang is set ...
		if b.ShellBang {
			b := path.Base(pkg)
			if err := af.AddRecord(cpio.StaticFile(filepath.Join(opts.BinaryDir, b), "#!/bbin/bb #!"+b+"\n", 0o755)); err != nil {
				return err
			}
		} else if err := af.AddRecord(cpio.Symlink(filepath.Join(opts.BinaryDir, path.Base(pkg)), "bb")); err != nil {
			return err
		}
	}
	return nil
}
