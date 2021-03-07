// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

Command gosimports updates your Go import lines,
adding missing ones, removing unreferenced ones,
and removing redundant blank lines.

     $ go get github.com/rinchsan/gosimports/cmd/gosimports

In addition to fixing imports, gosimports also formats
your code in the same style as gofmt so it can be used
as a replacement for your editor's gofmt-on-save hook.

For emacs, make sure you have the latest go-mode.el:
   https://github.com/dominikh/go-mode.el
Then in your .emacs file:
   (setq gofmt-command "gosimports")
   (add-hook 'before-save-hook 'gofmt-before-save)

For vim, set "gofmt_command" to "gosimports":
    https://golang.org/change/39c724dd7f252
    https://golang.org/wiki/IDEsAndTextEditorPlugins
    etc

For other editors, you probably know what to do.

To exclude directories in your $GOPATH from being scanned for Go
files, gosimports respects a configuration file at
$GOPATH/src/.gosimportsignore which may contain blank lines, comment
lines (beginning with '#'), or lines naming a directory relative to
the configuration file to ignore when scanning. No globbing or regex
patterns are allowed. Use the "-v" verbose flag to verify it's
working and see what gosimports is doing.

File bugs or feature requests at:

    https://github.com/rinchsan/gosimports/issues/new

Happy hacking!

*/
package main // import "github.com/rinchsan/gosimports/cmd/gosimports"
