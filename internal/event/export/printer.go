// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package export

import (
	"io"

	"github.com/rinchsan/gosimports/internal/event/core"
	"github.com/rinchsan/gosimports/internal/event/keys"
	"github.com/rinchsan/gosimports/internal/event/label"
)

type Printer struct {
	buffer [128]byte
}

func (p *Printer) WriteEvent(w io.Writer, ev core.Event, lm label.Map) {
	buf := p.buffer[:0]
	if !ev.At().IsZero() {
		_, _ = w.Write(ev.At().AppendFormat(buf, "2006/01/02 15:04:05 "))
	}
	msg := keys.Msg.Get(lm)
	_, _ = io.WriteString(w, msg)
	if err := keys.Err.Get(lm); err != nil {
		if msg != "" {
			_, _ = io.WriteString(w, ": ")
		}
		_, _ = io.WriteString(w, err.Error())
	}
	for index := 0; ev.Valid(index); index++ {
		l := ev.Label(index)
		if !l.Valid() || l.Key() == keys.Msg || l.Key() == keys.Err {
			continue
		}
		_, _ = io.WriteString(w, "\n\t")
		_, _ = io.WriteString(w, l.Key().Name())
		_, _ = io.WriteString(w, "=")
		l.Key().Format(w, buf, l)
	}
	_, _ = io.WriteString(w, "\n")
}
