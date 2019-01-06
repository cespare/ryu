// Copyright 2018 Ulf Adams
// Modifications copyright 2019 Caleb Spare
//
// The contents of this file may be used under the terms of the Apache License,
// Version 2.0.
//
//    (See accompanying file LICENSE or copy at
//     http://www.apache.org/licenses/LICENSE-2.0)
//
// Unless required by applicable law or agreed to in writing, this software
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.
//
// The code in this file is part of a Go translation of the C code written by
// Ulf Adams which may be found at https://github.com/ulfjack/ryu. That source
// code is licensed under Apache 2.0 and this code is derivative work thereof.

package ryu

import (
	"strconv"
	"testing"
)

func TestFormatFloat32(t *testing.T) {
	for _, f := range []float32{
		0,
		1,
		0.3,
		1000000,
		-123.45,
	} {
		got := FormatFloat32(f)
		want := strconv.FormatFloat(float64(f), 'e', -1, 32)
		if got != want {
			t.Errorf("FormatFloat32(%g): got %q; want %q", f, got, want)
		}
	}
}

var benchCases32 = []float32{
	0,
	1,
	0.3,
	1000000,
	-123.45,
}

var sink string

func BenchmarkFormatFloat32(b *testing.B) {
	for _, f := range benchCases32 {
		b.Run(FormatFloat32(f), func(b *testing.B) {
			var s string
			for i := 0; i < b.N; i++ {
				s = FormatFloat32(f)
			}
			sink = s
		})
	}
}

func BenchmarkStrconvFormatFloat32(b *testing.B) {
	for _, f := range benchCases32 {
		b.Run(FormatFloat32(f), func(b *testing.B) {
			var s string
			for i := 0; i < b.N; i++ {
				s = strconv.FormatFloat(float64(f), 'e', -1, 32)
			}
			sink = s
		})
	}
}
