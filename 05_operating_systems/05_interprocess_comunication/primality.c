#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
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
  bool (*func)(long), tty;

  if (argc != 2)
    exit_with_usage();

  if (strcmp(argv[1], "brute_force") == 0)
    func = &brute_force_impl;
  else if (strcmp(argv[1], "brutish") == 0)
    func = &brutish_impl;
  else if (strcmp(argv[1], "miller_rabin") == 0)
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
      int result = (*func)(num);
      fprintf(stdout, "%d\n", result);
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
      return 0;
  return 1;
}

int randint(int a, int b) { return rand() % (++b - a) + a; }

int modpow(int a, int d, int m)
{
  int c = a;
  for (int i = 1; i < d; i++)
    c = (c * a) % m;
  return c % m;
}

int witness(int a, int s, int d, int n)
{
  int x = modpow(a, d, n);
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
