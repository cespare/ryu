# ryu

[![GoDoc](https://godoc.org/github.com/cespare/ryu?status.svg)](https://godoc.org/github.com/cespare/ryu)

This is a Go implementation of [Ryu](https://github.com/ulfjack/ryu), a fast
algorithm for converting floating-point numbers to strings.

The API is:

```
func AppendFloat32(b []byte, f float32) []byte
func AppendFloat64(b []byte, f float64) []byte
func FormatFloat32(f float32) string
func FormatFloat64(f float64) string
```

These functions are the equivalents of calling strconv.FormatFloat or
strconv.AppendFloat using the formatter `'e'` and precision `-1`:

```
// These are the same:
const f float32 = 1.234
s := ryu.FormatFloat32(f)
s := strconv.FormatFloat(float64(f), 'e', -1, 32)
```

## Benchmarks

These benchmarks were taken with Go 1.12beta1 on Linux/amd64 using an
Intel i7-8700K.

```
name                                     old time/op    new time/op    delta
FormatFloat32-12                            128ns ± 1%      50ns ± 2%  -60.82%  (p=0.000 n=7+8)
FormatFloat64-12                            129ns ± 4%      65ns ± 5%  -49.54%  (p=0.000 n=7+8)
AppendFloat32/0e+00-12                     24.4ns ± 1%     3.0ns ± 1%  -87.88%  (p=0.000 n=8+8)
AppendFloat32/1e+00-12                     26.5ns ± 1%    13.2ns ± 3%  -49.98%  (p=0.000 n=8+8)
AppendFloat32/3e-01-12                     52.2ns ± 1%    32.5ns ± 2%  -37.73%  (p=0.000 n=8+8)
AppendFloat32/1e+06-12                     41.2ns ± 1%    17.9ns ± 1%  -56.45%  (p=0.000 n=8+7)
AppendFloat32/-1.2345e+02-12               83.3ns ± 2%    34.2ns ± 1%  -58.90%  (p=0.000 n=8+8)
AppendFloat64/0e+00-12                     24.5ns ± 2%     3.3ns ± 2%  -86.50%  (p=0.000 n=8+8)
AppendFloat64/1e+00-12                     26.9ns ± 1%    14.5ns ± 1%  -46.06%  (p=0.001 n=8+6)
AppendFloat64/3e-01-12                     53.0ns ± 1%    42.5ns ± 0%  -19.75%  (p=0.001 n=8+6)
AppendFloat64/1e+06-12                     41.4ns ± 1%    21.1ns ± 1%  -49.05%  (p=0.000 n=8+8)
AppendFloat64/-1.2345e+02-12               83.8ns ± 1%    43.3ns ± 1%  -48.32%  (p=0.000 n=8+8)
AppendFloat64/6.226662346353213e-309-12    25.5µs ± 1%     0.0µs ± 1%  -99.84%  (p=0.000 n=8+8)
```

The test `TestRandomBenchmark` gathers statistics about the distribution of call
latencies for random float64 values. Here is the summary for one sample of 10,000
random floats:

```
    ryu_test.go:279: after sampling 50000 float64s:
        ryu:               min = 2ns  max = 90ns     median = 41ns   mean = 41ns
        strconv (stdlib):  min = 8ns  max = 25845ns  median = 106ns  mean = 154ns
```

The `strconv.FormatFloat` latency is bimodal because of an infrequently-taken
slow path that is orders of magnitude more expensive
(https://golang.org/issue/15672).

## Size optimization

The Ryu algorithm requires several lookup tables. Ulf Adams's C library
implements a size optimization (`RYU_OPTIMIZE_SIZE`) which greatly reduces the
size of the float64 tables in exchange for a little more CPU cost.

I have a WIP implementation of this optimization on the `size` branch. A binary
built using that version is 7.96 kB smaller. The benchmark results take a hit as
compared with the non-size-optimized build:

```
name                                     old time/op    new time/op    delta
FormatFloat32-12                           50.0ns ± 2%    49.4ns ± 1%     ~     (p=0.183 n=8+8)
FormatFloat64-12                           65.0ns ± 5%    72.1ns ± 5%  +10.96%  (p=0.000 n=8+8)
AppendFloat32/0e+00-12                     2.95ns ± 1%    2.98ns ± 1%     ~     (p=0.072 n=8+8)
AppendFloat32/1e+00-12                     13.2ns ± 3%    13.1ns ± 1%     ~     (p=0.275 n=8+8)
AppendFloat32/3e-01-12                     32.5ns ± 2%    32.4ns ± 1%     ~     (p=0.742 n=8+8)
AppendFloat32/1e+06-12                     17.9ns ± 1%    17.6ns ± 1%   -2.12%  (p=0.001 n=7+8)
AppendFloat32/-1.2345e+02-12               34.2ns ± 1%    34.4ns ± 1%     ~     (p=0.426 n=8+8)
AppendFloat64/0e+00-12                     3.31ns ± 2%    3.29ns ± 1%     ~     (p=0.394 n=8+8)
AppendFloat64/1e+00-12                     14.5ns ± 1%    14.6ns ± 4%     ~     (p=0.641 n=6+8)
AppendFloat64/3e-01-12                     42.5ns ± 0%    50.0ns ± 1%  +17.44%  (p=0.001 n=6+8)
AppendFloat64/1e+06-12                     21.1ns ± 1%    21.1ns ± 2%     ~     (p=0.452 n=8+8)
AppendFloat64/-1.2345e+02-12               43.3ns ± 1%    50.9ns ± 1%  +17.57%  (p=0.000 n=8+8)
AppendFloat64/6.226662346353213e-309-12    40.6ns ± 1%    47.7ns ± 1%  +17.38%  (p=0.000 n=8+8)
```

However, it's still generally faster than strconv:

```
name                                     old time/op    new time/op    delta
FormatFloat32-12                            129ns ± 2%      49ns ± 1%  -61.72%  (p=0.000 n=8+8)
FormatFloat64-12                            130ns ± 3%      72ns ± 5%  -44.32%  (p=0.000 n=7+8)
AppendFloat32/0e+00-12                     24.5ns ± 2%     3.0ns ± 1%  -87.83%  (p=0.000 n=8+8)
AppendFloat32/1e+00-12                     26.4ns ± 1%    13.1ns ± 1%  -50.26%  (p=0.000 n=7+8)
AppendFloat32/3e-01-12                     52.6ns ± 2%    32.4ns ± 1%  -38.43%  (p=0.000 n=8+8)
AppendFloat32/1e+06-12                     41.3ns ± 2%    17.6ns ± 1%  -57.51%  (p=0.000 n=8+8)
AppendFloat32/-1.2345e+02-12               83.5ns ± 1%    34.4ns ± 1%  -58.82%  (p=0.000 n=8+8)
AppendFloat64/0e+00-12                     24.6ns ± 2%     3.3ns ± 1%  -86.63%  (p=0.000 n=8+8)
AppendFloat64/1e+00-12                     26.7ns ± 1%    14.6ns ± 4%  -45.51%  (p=0.000 n=8+8)
AppendFloat64/3e-01-12                     52.7ns ± 1%    50.0ns ± 1%   -5.17%  (p=0.000 n=8+8)
AppendFloat64/1e+06-12                     41.2ns ± 1%    21.1ns ± 2%  -48.61%  (p=0.000 n=7+8)
AppendFloat64/-1.2345e+02-12               83.7ns ± 1%    50.9ns ± 1%  -39.17%  (p=0.000 n=8+8)
AppendFloat64/6.226662346353213e-309-12    25.8µs ± 2%     0.0µs ± 1%  -99.81%  (p=0.000 n=8+8)
```

## Notes

This package is a fairly direct Go translation of Ulf Adams's C library at
https://github.com/ulfjack/ryu. This code is also licensed with Apache 2.0 as a
derived work of that code.

This package requires Go 1.12 (expected to be released February 2019).

For a small fraction of inputs, Ryu gives a different value than strconv does
for the last digit. This is due to a bug in strconv: https://golang.org/issue/29491.

## Future work

My plan is to incorporate this into strconv (see
https://golang.org/issue/15672). Then everyone will benefit from the faster
algorithm and there will be no need for this library.

If you would like to contribute, I'm interested in any bugfixes or clear-cut
optimizations, but given the above I don't intend to add more features or APIs
to this package.
