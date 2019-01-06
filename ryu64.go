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
	"math/bits"
)

const (
	mantBits64 = 52
	expBits64  = 11
	bias64     = 1023
)

func decimalLen32(u uint32) int {
	// Function precondition: u is not a 10-digit number.
	// (9 digits are sufficient for round-tripping.)
	assert(u < 1000000000, "too big")
	switch {
	case u >= 100000000:
		return 9
	case u >= 10000000:
		return 8
	case u >= 1000000:
		return 7
	case u >= 100000:
		return 6
	case u >= 10000:
		return 5
	case u >= 1000:
		return 4
	case u >= 100:
		return 3
	case u >= 10:
		return 2
	default:
		return 1
	}

}

// log10Pow2 returns floor(log_10(2^e)).
func log10Pow2(e int32) uint32 {
	// The first value this approximation fails for is 2^1651
	// which is just greater than 10^297.
	assert(e >= 0, "e >= 0")
	assert(e <= 1650, "e <= 1650")
	return (uint32(e) * 78913) >> 18
}

// log10Pow5 returns floor(log_10(5^e)).
func log10Pow5(e int32) uint32 {
	// The first value this approximation fails for is 5^2621
	// which is just greater than 10^1832.
	assert(e >= 0, "e >= 0")
	assert(e <= 2620, "e <= 2620")
	return (uint32(e) * 732923) >> 20
}

// pow5Bits returns ceil(log_2(5^e)), or else 1 if e==0.
func pow5Bits(e int32) int32 {
	// This approximation works up to the point that the multiplication
	// overflows at e = 3529. If the multiplication were done in 64 bits,
	// it would fail at 5^4004 which is just greater than 2^9297.
	assert(e >= 0, "e >= 0")
	assert(e <= 3528, "e <= 3528")
	return int32((uint32(e)*1217359)>>19 + 1)
}

func mulShift(m uint32, factor uint64, shift int32) uint32 {
	assert(shift > 32, "shift > 32")

	factorLo := uint32(factor)
	factorHi := uint32(factor >> 32)
	bits0 := uint64(m) * uint64(factorLo)
	bits1 := uint64(m) * uint64(factorHi)

	sum := (bits0 >> 32) + bits1
	shiftedSum := sum >> uint(shift-32)
	assert(shiftedSum <= math.MaxUint32, "shiftedSum <= math.MaxUint32")
	return uint32(shiftedSum)
}

func mulPow5InvDivPow2(m, q uint32, j int32) uint32 {
	return mulShift(m, pow5InvSplit[q], j)
}

func mulPow5DivPow2(m, i uint32, j int32) uint32 {
	return mulShift(m, pow5Split[i], j)
}

// FIXME: use bits.Div

func pow5Factor(v uint32) uint32 {
	for n := uint32(0); ; n++ {
		q, r := v/5, v%5
		if r != 0 {
			return n
		}
		v = q
	}
}

// multipleOfPowerOf5 reports whether v is divisible by 5^p.
func multipleOfPowerOf5(v, p uint32) bool {
	return pow5Factor(v) >= p
}

// multipleOfPowerOf5 reports whether v is divisible by 2^p.
func multipleOfPowerOf2(v, p uint32) bool {
	return uint32(bits.TrailingZeros32(v)) >= p
}
