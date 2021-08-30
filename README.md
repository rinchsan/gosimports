**CLI interface (such as options) is likely to change a few times before we reach 1.0.0**

# gosimports - simpler goimports

![](https://github.com/rinchsan/gosimports/workflows/CI/badge.svg)
![](https://img.shields.io/github/release/rinchsan/gosimports.svg?colorB=7E7E7E)
[![](https://pkg.go.dev/badge/github.com/rinchsan/gosimports.svg)](https://pkg.go.dev/github.com/rinchsan/gosimports/cmd/gosimports)

- :rocket: Drop-in replacement for `goimports`
- :100: Prettier than `goimports`
- :hammer: Originally forked from `golang.org/x/tools/cmd/goimports`

## Motivation

This `gosimports` provides one solution to the [goimports grouping/ordering problem](https://github.com/golang/go/issues/20818).

## Installation

### Go 1.16 or later

```bash
go install github.com/rinchsan/gosimports/cmd/gosimports@latest
```

### Go 1.15 or earlier

```bash
go get github.com/rinchsan/gosimports/cmd/gosimports
```

### Homebrew

```bash
brew install rinchsan/tap/gosimports
```

### Binary

Download binaries from [GitHub Releases](https://github.com/rinchsan/gosimports/releases)

## Example

```go
import (
	"bufio"

	// basic comments

	/*
		block comments
	*/

	"github.com/rinchsan/gosimports/internal/imports"

	"errors"
	gocmd "github.com/rinchsan/gosimports/internal/gocommand"
	"flag"

	"runtime"
	_ "runtime/pprof" // trailing inline comments
	"strings"

)
```

↓ `$ gosimports -w` :+1:

```go
import (
	"bufio"
	"errors"
	"flag"
	"runtime"
	_ "runtime/pprof" // trailing inline comments
	"strings"

	gocmd "github.com/rinchsan/gosimports/internal/gocommand"
	"github.com/rinchsan/gosimports/internal/imports"
)
```

while `goimports` formatting like below :-1:

```go
import (
	"bufio"

	// basic comments

	/*
		block comments
	*/

	"github.com/rinchsan/gosimports/internal/imports"

	"errors"
	"flag"

	gocmd "github.com/rinchsan/gosimports/internal/gocommand"

	"runtime"
	_ "runtime/pprof" // trailing inline comments
	"strings"
)
```

## License

Copyright 2013 The Go Authors. All rights reserved.

Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.
