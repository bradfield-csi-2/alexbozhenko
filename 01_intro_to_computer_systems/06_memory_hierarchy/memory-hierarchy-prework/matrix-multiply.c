/*
Naive code for multiplying two matrices together.

There must be a better way!
*/

#include <stdio.h>
#include <stdlib.h>

/*
  A naive implementation of matrix multiplication.

  DO NOT MODIFY THIS FUNCTION, the tests assume it works correctly, which it
  currently does
*/
void matrix_multiply(double **C, double **A, double **B, int a_rows, int a_cols,
                     int b_cols)
{
  for (int i = 0; i < a_rows; i++)
  {
    for (int j = 0; j < b_cols; j++)
    {
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++)
        C[i][j] += A[i][k] * B[k][j];
    }
  }
}

void fast_matrix_multiply(double **c, double **a, double **b, int a_rows,
                          int a_cols, int b_cols)
{

  // return matrix_multiply(c, a, b, a_rows, a_cols, b_cols);

  // for (int a_row = 0; a_row < a_rows; a_row++)
  // {
  //   for (int b_col = 0; b_col < b_cols; b_col++)
  //   {
  //     double c_cell = 0;
  //     //c[a_row][b_col] = 0;
  //     for (int a_col = 0; a_col < a_cols; a_col++)
  //     {
  //       c_cell += a[a_row][a_col] * b[a_col][b_col];
  //     }
  //     c[a_row][b_col] = c_cell;
  //   }
  // }

  //  int b_rows = a_cols;

  for (int a_row = 0; a_row < a_rows; a_row++)
  {
    for (int a_col = 0; a_col < a_cols; a_col++)
    {
      for (int b_col = 0; b_col < b_cols; b_col++)
      {
        if (a_col == 0)
        {
          // printf("my na a_col=%d\tzanylayem! c[a_row=%d][b_col=%d]\n", a_col, a_row, b_col);
          c[a_row][b_col] = 0;
        }
        c[a_row][b_col] += a[a_row][a_col] * b[a_col][b_col];
      }
    }
  }
}
