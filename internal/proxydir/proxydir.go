// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package proxydir provides functions for writing module data to a directory
// in proxy format, so that it can be used as a module proxy by setting
// GOPROXY="file://<dir>".
package proxydir

import (
	"path/filepath"
	"strings"

	"github.com/rinchsan/gosimports/internal/testenv"
)

// ToURL returns the file uri for a proxy directory.
func ToURL(dir string) string {
	if testenv.Go1Point() >= 13 {
		// file URLs on Windows must start with file:///. See golang.org/issue/6027.
		path := filepath.ToSlash(dir)
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		return "file://" + path
	} else {
		// Prior to go1.13, the Go command on Windows only accepted GOPROXY file URLs
		// of the form file://C:/path/to/proxy. This was incorrect: when parsed, "C:"
		// is interpreted as the host. See golang.org/issue/6027. This has been
		// fixed in go1.13, but we emit the old format for old releases.
		return "file://" + filepath.ToSlash(dir)
	}
}
