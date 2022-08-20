// Package gosimports implements a Go pretty-printer
// that updates your Go import lines, adding missing ones, removing unreferenced ones,
// and removing redundant blank lines.
package gosimports

import (
	"log"
	"os"

	"github.com/rinchsan/gosimports/internal/gocommand"
	"github.com/rinchsan/gosimports/internal/imports"
)

// Options specifies options for processing files.
type Options struct {
	Fragment  bool // Accept fragment of a source file (no package statement)
	AllErrors bool // Report all errors (not just the first 10 on different lines)

	Comments  bool // Print comments (true if nil *Options provided)
	TabIndent bool // Use tabs for indent (true if nil *Options provided)
	TabWidth  int  // Tab width (8 if nil *Options provided)

	FormatOnly bool // Disable the insertion and deletion of imports
}

// Debug controls verbose logging.
var Debug = false

// LocalPrefix is a comma-separated string of import path prefixes, which, if
// set, instructs Process to sort the import paths with the given prefixes
// into another group after 3rd-party packages.
var LocalPrefix string

// Process formats and adjusts imports for the provided file.
// If opt is nil, the defaults are used.
// If src is nil, the source is read from the filesystem.
//
// Note that filename's directory influences which imports can be chosen,
// so it is important that filename be accurate.
// To process data “as if” it were in filename, pass the data as a non-nil src.
func Process(filename string, src []byte, opt *Options) ([]byte, error) {
	var err error
	if src == nil {
		src, err = os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	}
	if opt == nil {
		opt = &Options{
			Fragment:  true,
			Comments:  true,
			TabIndent: true,
			TabWidth:  8,
		}
	}
	intopt := &imports.Options{
		Env: &imports.ProcessEnv{
			GocmdRunner: &gocommand.Runner{},
		},
		LocalPrefix: LocalPrefix,
		Fragment:    opt.Fragment,
		AllErrors:   opt.AllErrors,
		Comments:    opt.Comments,
		TabIndent:   opt.TabIndent,
		TabWidth:    opt.TabWidth,
		FormatOnly:  opt.FormatOnly,
	}
	if Debug {
		intopt.Env.Logf = log.Printf
	}
	return imports.Process(filename, src, intopt)
}
