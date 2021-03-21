# gosimports - simpler goimports

![](https://github.com/rinchsan/gosimports/workflows/CI/badge.svg)
![](https://img.shields.io/github/release/rinchsan/gosimports.svg?colorB=7E7E7E)
[![](https://pkg.go.dev/badge/github.com/rinchsan/gosimports.svg)](https://pkg.go.dev/github.com/rinchsan/gosimports/cmd/gosimports)

- :rocket: Just replace `goimports` with `gosimports` !
- :sparkles: Simpler than `goimports`
- :100: Interface compatible with `goimports` (result is not compatible, but simpler)
- :hammer: Originally forked from `golang.org/x/tools/cmd/goimports`

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

	// comment

	"golang.org/x/tools/internal/imports"

	"errors"
	gocmd "golang.org/x/tools/internal/gocommand"
	"flag"

	"runtime" // comment
	_ "runtime/pprof"
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
	_ "runtime/pprof"
	"strings"

	gocmd "golang.org/x/tools/internal/gocommand"
	"golang.org/x/tools/internal/imports"
)
```

while `goimports` formatting like below :-1:

```go
import (
	"bufio"

	// comment

	"golang.org/x/tools/internal/imports"

	"errors"
	"flag"

	gocmd "golang.org/x/tools/internal/gocommand"

	"runtime" // comment
	_ "runtime/pprof"
	"strings"
)
```

## License

Copyright 2013 The Go Authors. All rights reserved.

Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.
