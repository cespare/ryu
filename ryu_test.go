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
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
)

var genericTestCases = []float64{
	0,
	math.Float64frombits(uint64(1) << 63), // -0
	math.NaN(),
	math.Inf(-1),
	math.Inf(1),
	1,
	-1,
	10,
	-10,
	0.3,
	-0.3,
	1000000,
	123456.7,
	123e45,
	-123.45,
	1e23,
	math.SmallestNonzeroFloat32,
	math.MaxFloat32,
	below1e23,
	above1e23,

	// https://golang.org/issue/2625
	383260575764816448,
	383260575764816448,
}

const (
	below1e23 = 99999999999999974834176
	above1e23 = 100000000000000008388608
)

func TestFormatFloat32(t *testing.T) {
	for _, f64 := range genericTestCases {
		f := float32(f64)
		got := FormatFloat32(f)
		want := strconv.FormatFloat(float64(f), 'e', -1, 32)
		if got != want {
			t.Errorf("FormatFloat32(%g): got %q; want %q", f, got, want)
		}
	}
}

var float64TestCases = []float64{
	123e300,
	123e-300,
	5e-324,
	-5e-324,
	math.SmallestNonzeroFloat64,
	math.MaxFloat64,

	// https://www.exploringbinary.com/java-hangs-when-converting-2-2250738585072012e-308/
	2.2250738585072012e-308,
	// https://www.exploringbinary.com/php-hangs-on-numeric-value-2-2250738585072011e-308/
	2.2250738585072011e-308,

	// https://github.com/golang/go/issues/29491
	//498484681984085570,
	//-5.8339553793802237e+23,
}

func TestFormatFloat64(t *testing.T) {
	for _, f := range append(genericTestCases, float64TestCases...) {
		got := FormatFloat64(f)
		want := strconv.FormatFloat(f, 'e', -1, 64)
		if got != want {
			t.Errorf("FormatFloat64(%g): got %q; want %q", f, got, want)
		}
	}
}

func TestFormatFloatRandom(t *testing.T) {
	t.Skip("disabled because of Go bug: https://github.com/golang/go/issues/29491")
	for i := 0; i < 1e6; i++ {
		f := math.Float64frombits(rand.Uint64())

		got32 := FormatFloat32(float32(f))
		want32 := strconv.FormatFloat(f, 'e', -1, 32)
		if got32 != want32 {
			t.Fatalf("FormatFloat32(%g): got %q; want %q", f, got32, want32)
		}

		got := FormatFloat64(f)
		want := strconv.FormatFloat(f, 'e', -1, 64)
		if got != want {
			t.Fatalf("FormatFloat64(%g): got %q; want %q", f, got, want)
		}
	}
}

var sink string
var sinkb []byte

// Much of the Format cost is allocation, so most of the interesting benchmarks
// are for Append, below.

const benchFloat = 123.456

func BenchmarkFormatFloat32(b *testing.B) {
	var s string
	f := float32(benchFloat)
	for i := 0; i < b.N; i++ {
		s = FormatFloat32(f)
	}
	sink = s
}

func BenchmarkStrconvFormatFloat32(b *testing.B) {
	var s string
	f := float32(benchFloat)
	for i := 0; i < b.N; i++ {
		s = strconv.FormatFloat(float64(f), 'e', -1, 32)
	}
	sink = s
}

func BenchmarkFormatFloat64(b *testing.B) {
	var s string
	f := float64(benchFloat)
	for i := 0; i < b.N; i++ {
		s = FormatFloat64(f)
	}
	sink = s
}

func BenchmarkStrconvFormatFloat64(b *testing.B) {
	var s string
	f := float64(benchFloat)
	for i := 0; i < b.N; i++ {
		s = strconv.FormatFloat(f, 'e', -1, 64)
	}
	sink = s
}

var benchCases = []float64{
	0,
	1,
	0.3,
	1000000,
	-123.45,
}

func BenchmarkAppendFloat32(b *testing.B) {
	for _, f64 := range benchCases {
		f := float32(f64)
		b.Run(FormatFloat32(f), func(b *testing.B) {
			var buf []byte
			for i := 0; i < b.N; i++ {
				buf = AppendFloat32(buf[:0], f)
			}
			sinkb = buf
		})
	}
}

func BenchmarkStrconvAppendFloat32(b *testing.B) {
	for _, f64 := range benchCases {
		f := float32(f64)
		b.Run(FormatFloat32(f), func(b *testing.B) {
			var buf []byte
			for i := 0; i < b.N; i++ {
				buf = strconv.AppendFloat(buf[:0], float64(f), 'e', -1, 32)
			}
			sinkb = buf
		})
	}
}

var benchCases64 = []float64{
	622666234635.3213e-320, // https://golang.org/issue/15672
}

func BenchmarkAppendFloat64(b *testing.B) {
	for _, f := range append(benchCases, benchCases64...) {
		b.Run(FormatFloat64(f), func(b *testing.B) {
			var buf []byte
			for i := 0; i < b.N; i++ {
				buf = AppendFloat64(buf[:0], f)
			}
			sinkb = buf
		})
	}
}

func BenchmarkStrconvAppendFloat64(b *testing.B) {
	for _, f := range append(benchCases, benchCases64...) {
		b.Run(FormatFloat64(f), func(b *testing.B) {
			var buf []byte
			for i := 0; i < b.N; i++ {
				buf = strconv.AppendFloat(buf[:0], f, 'e', -1, 64)
			}
			sinkb = buf
		})
	}
}

func TestDecimalLen(t *testing.T) {
	for n := uint64(1); n < 1000; n++ {
		testDecimalLen(t, n)
	}
	for i := 0; i < 1e5; i++ {
		n := uint64(rand.Intn(99999999999999999) + 1)
		testDecimalLen(t, n)
	}
}

func testDecimalLen(t *testing.T, n uint64) {
	t.Helper()
	want := len(big.NewInt(int64(n)).String()) // n fits into int64
	if got := decimalLen64(n); got != want {
		t.Fatalf("decimalLen64(%d): got %d; want %d", n, got, want)
	}
	if n < math.MaxUint32 {
		if got := decimalLen32(uint32(n)); got != want {
			t.Fatalf("decimalLen32(%d): got %d; want %d", n, got, want)
		}
	}
}
