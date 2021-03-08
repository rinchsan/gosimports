// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metric

import (
	"github.com/rinchsan/gosimports/internal/event/keys"
	"github.com/rinchsan/gosimports/internal/event/label"
)

// Scalar represents the construction information for a scalar metric.
type Scalar struct {
	// Name is the unique name of this metric.
	Name string
	// Description can be used by observers to describe the metric to users.
	Description string
	// Keys is the set of labels that collectively describe rows of the metric.
	Keys []label.Key
}

// HistogramInt64 represents the construction information for an int64 histogram metric.
type HistogramInt64 struct {
	// Name is the unique name of this metric.
	Name string
	// Description can be used by observers to describe the metric to users.
	Description string
	// Keys is the set of labels that collectively describe rows of the metric.
	Keys []label.Key
	// Buckets holds the inclusive upper bound of each bucket in the histogram.
	Buckets []int64
}

// HistogramFloat64 represents the construction information for an float64 histogram metric.
type HistogramFloat64 struct {
	// Name is the unique name of this metric.
	Name string
	// Description can be used by observers to describe the metric to users.
	Description string
	// Keys is the set of labels that collectively describe rows of the metric.
	Keys []label.Key
	// Buckets holds the inclusive upper bound of each bucket in the histogram.
	Buckets []float64
}

// SumInt64 creates a new metric based on the Scalar information that sums all
// the values recorded on the int64 measure.
// Metrics of this type will use Int64Data.
func (info Scalar) SumInt64(e *Config, key *keys.Int64) {
	data := &Int64Data{Info: &info, key: key}
	e.subscribe(key, data.sum)
}

// Record creates a new metric based on the HistogramInt64 information that
// tracks the bucketized counts of values recorded on the int64 measure.
// Metrics of this type will use HistogramInt64Data.
func (info HistogramInt64) Record(e *Config, key *keys.Int64) {
	data := &HistogramInt64Data{Info: &info, key: key}
	e.subscribe(key, data.record)
}

// Record creates a new metric based on the HistogramFloat64 information that
// tracks the bucketized counts of values recorded on the float64 measure.
// Metrics of this type will use HistogramFloat64Data.
func (info HistogramFloat64) Record(e *Config, key *keys.Float64) {
	data := &HistogramFloat64Data{Info: &info, key: key}
	e.subscribe(key, data.record)
}
