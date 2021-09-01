#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <mqueue.h>
#include <time.h>
#include <bsd/sys/time.h>
#include <errno.h>
#include "types.h"

bool brute_force_impl(long n);
bool brutish_impl(long n);
bool miller_rabin_impl(long n);

__attribute__((__noreturn__)) void exit_with_usage(void)
{
  fprintf(stderr, "Usage: ./primality [brute_force|brutish|miller_rabin]\n");
  exit(1);
}

int main(int argc, char *argv[])
{
  long num;
  bool (*func)(long), tty, result;

  if (argc == 2) // console mode
  {
    if (strcmp(argv[1], ALGORITHMS_STRING[brute_force]) == 0)
      func = &brute_force_impl;
    else if (strcmp(argv[1], ALGORITHMS_STRING[brutish]) == 0)
      func = &brutish_impl;
    else if (strcmp(argv[1], ALGORITHMS_STRING[miller_rabin]) == 0)
      func = &miller_rabin_impl;
    else
      exit_with_usage();

    tty = isatty(fileno(stdin));

    if (tty)
    {
      fprintf(stderr, "Running \"%s\", enter a number:\n> ", argv[1]);

      while (scanf("%ld", &num) == 1)
      {
        printf("%d\n", (*func)(num));
        fflush(stdout);
        fprintf(stderr, "> ");
      }
    }
    else
    {
      int elements_read;

      while ((elements_read = scanf("%ld", &num)) != EOF && elements_read == 1)
      {
        result = (*func)(num);
        fprintf(stdout, "%d\n", result);
      }
    }
  }
  else // IPC worker mode
  {

    mqd_t request_queue = mq_open(request_queue_name, O_RDWR);

    if (request_queue == (mqd_t)-1)
    {
      printf("%s\n", strerror(errno));
      exit(EXIT_FAILURE);
    }
    mqd_t response_queue = mq_open(response_queue_name, O_RDWR);
    if (request_queue == (mqd_t)-1)
    {
      printf("%s\n", strerror(errno));
      exit(EXIT_FAILURE);
    }

    struct request req;
    struct response resp;

    while (true)
    {
      mq_receive(request_queue, (char *)&req,
                 sizeof(struct request), (unsigned int)0);
      //  printf("received request for %ld with alg %s\n", req.number,             ALGORITHMS_STRING[req.alg]);
      switch (req.alg)
      {
      case brute_force:
        func = &brute_force_impl;
        break;
      case brutish:
        func = &brutish_impl;
        break;
      case miller_rabin:
        func = &miller_rabin_impl;
        break;
      default:
        fprintf(stderr, "Got request for unexpected algorithm\n");
        exit(-1);
      }

      struct timespec begin, end, time_spent;
      clock_gettime(CLOCK_PROCESS_CPUTIME_ID, &begin);
      result = (*func)(req.number);
      clock_gettime(CLOCK_PROCESS_CPUTIME_ID, &end);
      timespecsub(&end, &begin, &time_spent);
      // printf("time_spent %lld.%.9ld\n", (long long)time_spent.tv_sec, time_spent.tv_nsec);
      double time_spent_secs = (double)(time_spent.tv_sec) +
                               (double)(time_spent.tv_nsec) / 1000000000;

      resp.alg = req.alg;
      resp.number = req.number;
      resp.result = result;
      resp.duration = time_spent_secs;
      mq_send(response_queue, (const char *)&resp,
              sizeof(struct response), (unsigned int)0);
    }
  }
}

/*
 * Primality test implementations
 */

// Just test every factor
bool brute_force_impl(long n)
{
  for (long i = 2; i < n; i++)
    if (n % i == 0)
      return 0;
  return 1;
}

// Test factors, up to sqrt(n)
bool brutish_impl(long n)
{
  long max = (long)floor(sqrt((double)n));
  for (long i = 2; i <= max; i++)
    if (n % i == 0)
    {
      return 0;
    }
  return 1;
}

int randint(int a, int b) { return rand() % (++b - a) + a; }

long modpow(long a, long d, long m)
{
  long c = a;
  for (int i = 1; i < d; i++)
    c = (c * a) % m;
  return c % m;
}

int witness(int a, int s, int d, int n)
{
  long x = modpow(a, d, n);
  if (x == 1)
    return 1;
  for (int i = 0; i < s - 1; i++)
  {
    if (x == n - 1)
      return 1;
    x = modpow(x, 2, n);
  }
  return (x == n - 1);
}

// TODO we should probably make this a parameter!
int MILLER_RABIN_ITERATIONS = 10;

// An implementation of the probabilistic Miller-Rabin test
bool miller_rabin_impl(long n)
{
  int a, s = 0, d = (int)n - 1;

  if (n == 2)
    return 1;

  if (!(n & 1) || n <= 1)
    return 0;

  while (!(d & 1))
  {
    d >>= 1;
    s += 1;
  }
  for (int i = 0; i < MILLER_RABIN_ITERATIONS; i++)
  {
    a = randint(2, (int)n - 1);
    if (!witness(a, s, d, (int)n))
      return 0;
  }
  return 1;
}
