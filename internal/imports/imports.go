// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run mkstdlib.go

// Package imports implements a Go pretty-printer (like package "go/format")
// that also adds or removes import statements as necessary.
package imports

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// Options is golang.org/x/tools/imports.Options with extra internal-only options.
type Options struct {
	Env *ProcessEnv // The environment to use. Note: this contains the cached module and filesystem state.

	// LocalPrefix is a comma-separated string of import path prefixes, which, if
	// set, instructs Process to sort the import paths with the given prefixes
	// into another group after 3rd-party packages.
	LocalPrefix string

	Fragment  bool // Accept fragment of a source file (no package statement)
	AllErrors bool // Report all errors (not just the first 10 on different lines)

	Comments  bool // Print comments (true if nil *Options provided)
	TabIndent bool // Use tabs for indent (true if nil *Options provided)
	TabWidth  int  // Tab width (8 if nil *Options provided)

	FormatOnly bool // Disable the insertion and deletion of imports
}

// Process implements golang.org/x/tools/imports.Process with explicit context in opt.Env.
func Process(filename string, src []byte, opt *Options) (formatted []byte, err error) {
	fileSet := token.NewFileSet()
	file, adjust, err := parse(fileSet, filename, src, opt)
	if err != nil {
		return nil, err
	}

	if !opt.FormatOnly {
		if err := fixImports(fileSet, file, filename, opt.Env); err != nil {
			return nil, err
		}
	}
	return formatFile(fileSet, file, src, adjust, opt)
}

// formatFile formats the file syntax tree.
// It may mutate the token.FileSet.
//
// If an adjust function is provided, it is called after formatting
// with the original source (formatFile's src parameter) and the
// formatted file, and returns the postpocessed result.
func formatFile(fset *token.FileSet, file *ast.File, src []byte, adjust func(orig []byte, src []byte) []byte, opt *Options) ([]byte, error) {
	mergeImports(file)
	sortImports(opt.LocalPrefix, fset.File(file.Pos()), file)
	impsByGroup := make(map[int][]*ast.ImportSpec)
	for _, impSection := range astutil.Imports(fset, file) {
		for _, importSpec := range impSection {
			importPath, _ := strconv.Unquote(importSpec.Path.Value)
			groupNum := importGroup(opt.LocalPrefix, importPath)
			impsByGroup[groupNum] = append(impsByGroup[groupNum], importSpec)
		}
	}

	printerMode := printer.UseSpaces
	if opt.TabIndent {
		printerMode |= printer.TabIndent
	}
	printConfig := &printer.Config{Mode: printerMode, Tabwidth: opt.TabWidth}

	var buf bytes.Buffer
	if err := printConfig.Fprint(&buf, fset, file); err != nil {
		return nil, err
	}
	out := buf.Bytes()
	if adjust != nil {
		out = adjust(src, out)
	}
	out, err := separateImportsIntoGroups(bytes.NewReader(out), impsByGroup)
	if err != nil {
		return nil, err
	}

	out, err = format.Source(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// parse parses src, which was read from filename,
// as a Go source file or statement list.
func parse(fset *token.FileSet, filename string, src []byte, opt *Options) (*ast.File, func(orig, src []byte) []byte, error) {
	parserMode := parser.Mode(0)
	if opt.Comments {
		parserMode |= parser.ParseComments
	}
	if opt.AllErrors {
		parserMode |= parser.AllErrors
	}

	// Try as whole source file.
	file, err := parser.ParseFile(fset, filename, src, parserMode)
	if err == nil {
		return file, nil, nil
	}
	// If the error is that the source file didn't begin with a
	// package line and we accept fragmented input, fall through to
	// try as a source fragment.  Stop and return on any other error.
	if !opt.Fragment || !strings.Contains(err.Error(), "expected 'package'") {
		return nil, nil, err
	}

	// If this is a declaration list, make it a source file
	// by inserting a package clause.
	// Insert using a ;, not a newline, so that parse errors are on
	// the correct line.
	const prefix = "package main;"
	psrc := append([]byte(prefix), src...)
	file, err = parser.ParseFile(fset, filename, psrc, parserMode)
	if err == nil {
		// Gofmt will turn the ; into a \n.
		// Do that ourselves now and update the file contents,
		// so that positions and line numbers are correct going forward.
		psrc[len(prefix)-1] = '\n'
		fset.File(file.Package).SetLinesForContent(psrc)

		// If a main function exists, we will assume this is a main
		// package and leave the file.
		if containsMainFunc(file) {
			return file, nil, nil
		}

		adjust := func(orig, src []byte) []byte {
			// Remove the package clause.
			src = src[len(prefix):]
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}
	// If the error is that the source file didn't begin with a
	// declaration, fall through to try as a statement list.
	// Stop and return on any other error.
	if !strings.Contains(err.Error(), "expected declaration") {
		return nil, nil, err
	}

	// If this is a statement list, make it a source file
	// by inserting a package clause and turning the list
	// into a function body.  This handles expressions too.
	// Insert using a ;, not a newline, so that the line numbers
	// in fsrc match the ones in src.
	fsrc := append(append([]byte("package p; func _() {"), src...), '}')
	file, err = parser.ParseFile(fset, filename, fsrc, parserMode)
	if err == nil {
		adjust := func(orig, src []byte) []byte {
			// Remove the wrapping.
			// Gofmt has turned the ; into a \n\n.
			src = src[len("package p\n\nfunc _() {"):]
			src = src[:len(src)-len("}\n")]
			// Gofmt has also indented the function body one level.
			// Remove that indent.
			src = bytes.Replace(src, []byte("\n\t"), []byte("\n"), -1)
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}

	// Failed, and out of options.
	return nil, nil, err
}

// containsMainFunc checks if a file contains a function declaration with the
// function signature 'func main()'
func containsMainFunc(file *ast.File) bool {
	for _, decl := range file.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			if f.Name.Name != "main" {
				continue
			}

			if len(f.Type.Params.List) != 0 {
				continue
			}

			if f.Type.Results != nil && len(f.Type.Results.List) != 0 {
				continue
			}

			return true
		}
	}

	return false
}

func cutSpace(b []byte) (before, middle, after []byte) {
	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\t' || b[i] == '\n') {
		i++
	}
	j := len(b)
	for j > 0 && (b[j-1] == ' ' || b[j-1] == '\t' || b[j-1] == '\n') {
		j--
	}
	if i <= j {
		return b[:i], b[i:j], b[j:]
	}
	return nil, nil, b[j:]
}

// matchSpace reformats src to use the same space context as orig.
//  1. If orig begins with blank lines, matchSpace inserts them at the beginning of src.
//  2. matchSpace copies the indentation of the first non-blank line in orig
//     to every non-blank line in src.
//  3. matchSpace copies the trailing space from orig and uses it in place
//     of src's trailing space.
func matchSpace(orig []byte, src []byte) []byte {
	before, _, after := cutSpace(orig)
	i := bytes.LastIndex(before, []byte{'\n'})
	before, indent := before[:i+1], before[i+1:]

	_, src, _ = cutSpace(src)

	var b bytes.Buffer
	b.Write(before)
	for len(src) > 0 {
		line := src
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line, src = line[:i+1], line[i+1:]
		} else {
			src = nil
		}
		if len(line) > 0 && line[0] != '\n' { // not blank
			b.Write(indent)
		}
		b.Write(line)
	}
	b.Write(after)
	return b.Bytes()
}

// separateImportsIntoGroups separates import lines into groups determined by importGroup func.
func separateImportsIntoGroups(r io.Reader, impsByGroup map[int][]*ast.ImportSpec) ([]byte, error) {
	var out bytes.Buffer
	in := bufio.NewReader(r)
	inImports := false
	done := false
	impInserted := false
	for {
		s, err := in.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if !inImports && !done && strings.HasPrefix(s, "import") && strings.Contains(s, "(") {
			inImports = true
			fmt.Fprint(&out, s)
			continue
		}
		if inImports && !impInserted {
			for i := 0; i <= 3; i++ {
				for _, imp := range impsByGroup[i] {
					if imp.Path.Value == `"C"` {
						continue
					}
					if imp.Name != nil {
						fmt.Fprint(&out, imp.Name.Name, " ")
					}
					fmt.Fprint(&out, imp.Path.Value)
					if imp.Comment != nil {
						for _, comment := range imp.Comment.List {
							fmt.Fprint(&out, " ", comment.Text)
						}
					}
					out.WriteByte('\n')
				}
				out.WriteByte('\n')
			}
			impInserted = true
			continue
		}
		if inImports && !strings.Contains(s, "//") && strings.Contains(s, ")") {
			done = true
			inImports = false
		}

		if !inImports {
			fmt.Fprint(&out, s)
		}
	}
	return out.Bytes(), nil
}
