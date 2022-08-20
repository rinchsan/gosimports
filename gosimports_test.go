package gosimports_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/rinchsan/gosimports"
)

func TestMain(m *testing.M) {
	gosimports.Debug = true
	m.Run()
}

func TestProcess_from_src(t *testing.T) {
	src, err := os.ReadFile("gosimports.go")
	if err != nil {
		t.Fatal("expected: err == nil")
	}

	formatted, err := gosimports.Process("", src, nil)
	if err != nil {
		t.Fatal("expected: err == nil")
	}

	if !bytes.Equal(src, formatted) {
		t.Fatal("expected: src == formatted")
	}
}

func TestProcess_from_filename(t *testing.T) {
	formatted, err := gosimports.Process("gosimports.go", nil, nil)
	if err != nil {
		t.Fatal("expected: err == nil")
	}

	src, err := os.ReadFile("gosimports.go")
	if err != nil {
		t.Fatal("expected: err == nil")
	}

	if !bytes.Equal(src, formatted) {
		t.Fatal("expected: src == formatted")
	}
}

func TestProcess_unknown_filename(t *testing.T) {
	_, err := gosimports.Process("unknown.go", nil, nil)
	if err == nil {
		t.Fatal("expected: err != nil")
	}
}
