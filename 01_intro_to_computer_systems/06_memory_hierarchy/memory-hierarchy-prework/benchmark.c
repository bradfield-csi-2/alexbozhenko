/*
  A very simple test and benchmark suite
*/
#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <time.h>
#include "matrix-multiply.h"

// Allocate space for an m x n matrix. Caller frees.
double **matrix_alloc(int m, int n)
{
  double **matrix = malloc(m * sizeof(double));

  for (int i = 0; i < m; i++)
    matrix[i] = calloc(n, sizeof(double));

  return matrix;
}

// Free the entirety of an m row matrix
void matrix_free(double **matrix, int m)
{
  for (int i = 0; i < m; i++)
    free(matrix[i]);
  free(matrix);
}

// Fill an m x n matrix with random values
void matrix_fill_random(double **matrix, int m, int n)
{
  for (int i = 0; i < m; i++)
    for (int j = 0; j < n; j++)
      matrix[i][j] = (double)rand() / (double)RAND_MAX;
}

// Verify that two m x n matrices contain the same values
bool matrix_equal(double **A, double **B, int m, int n)
{
  for (int i = 0; i < m; i++)
    for (int j = 0; j < n; j++)
      if (A[i][j] != B[i][j])
        return false;
  return true;
}

// To ensure fair cache pre-population, write junk to cache to "flush" it
void flush_cache()
{
  int size = 4 * 1024 * 1024; // 4MB to clear out L3
  char *b = malloc(size);
  for (int i = 0; i < size; i++)
    b[i] = i;
  free(b);
}

int main(int argc, char *argv[])
{
  if (argc != 2)
  {
    printf("Usage: ./benchmark [n]\n");
    exit(1);
  }

  int n = atoi(argv[1]);
  clock_t start, stop;
  double naive_time, fast_time;

  // alloc input and output matrices
  double **A = matrix_alloc(n, n);
  double **B = matrix_alloc(n, n);
  double **C_naive = matrix_alloc(n, n);
  double **C_fast = matrix_alloc(n, n);

  // input matrices should have random values
  matrix_fill_random(A, n, n);
  matrix_fill_random(B, n, n);

  // compute the product naively
  flush_cache();
  start = clock();
  matrix_multiply(C_naive, A, B, n, n, n);
  stop = clock();
  naive_time = (stop - start) / (double)CLOCKS_PER_SEC;

  // compute the product "quickly"
  flush_cache();
  start = clock();
  fast_matrix_multiply(C_fast, A, B, n, n, n);
  stop = clock();
  fast_time = (stop - start) / (double)CLOCKS_PER_SEC;

  printf("Naive: %.3fs\nFast: %.3fs\n%0.2fx speedup\n", naive_time, fast_time,
         naive_time / fast_time);

  // verify that both outputs are the same
  if (!matrix_equal(C_naive, C_fast, n, n))
  {
    printf("\nHowever, matrix results did not match!\n");
    for (int i = 0; i < n; i++)
    {
      for (int j = 0; j < n; j++)
      {
        printf("%f ", C_naive[i][j]);
      }
      printf("\n");
    }

    printf("\n");
    for (int i = 0; i < n; i++)
    {
      for (int j = 0; j < n; j++)
      {
        printf("%f ", C_fast[i][j]);
      }
      printf("\n");
    }
  }

  // free everything
  matrix_free(A, n);
  matrix_free(B, n);
  matrix_free(C_naive, n);
  matrix_free(C_fast, n);

  return 0;
}
