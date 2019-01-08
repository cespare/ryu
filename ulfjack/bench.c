#include <stdint.h>
#include <stdio.h>
#include <time.h>
#include <unistd.h>

#include "ryu/ryu.h"

int64_t time_sub(const struct timespec *t0, const struct timespec *t1) {
  int64_t nsec = (int64_t)t0->tv_sec * 1000000000 + (int64_t)t0->tv_nsec;
  nsec -= (int64_t)t1->tv_sec * 1000000000 + (int64_t)t1->tv_nsec;
  return nsec;
}

int main(int argc, char **argv) {
  struct timespec start, end;
  int64_t elapsed;
  int64_t iters = 0;

  char buf[40];
  int sink;

  clock_gettime(CLOCK_MONOTONIC, &start);
  for (;;) {
    for (int i = 0; i < 10000; i++) {
      d2s_buffered(1.0, buf);
      sink += buf[2];
    }
    clock_gettime(CLOCK_MONOTONIC, &end);

    iters += 10000;
    elapsed = time_sub(&end, &start);
    if (elapsed >= 1000000000) {
      break;
    }
  }

  double secs = (double)elapsed / 1000000000.0;
  printf("%lu iters in %lf secs: %.2lf ns/iter\n", iters, secs, (double)elapsed / (double)iters);
  if (argc == 1000) {
    printf("%d\n", sink);
  }
}
