/*

Two different ways to loop over an array of arrays.

Spotted at:
http://stackoverflow.com/questions/9936132/why-does-the-order-of-the-loops-affect-performance-when-iterating-over-a-2d-arra

*/
#define SIZE 550

void option_one()
{
  int i, j;
  static int x[SIZE][SIZE];
  for (i = 0; i < SIZE; i++)
  {
    for (j = 0; j < SIZE; j++)
    {
      x[i][j] = i + j;
    }
  }
}

void option_two()
{
  int i, j;
  static int x[SIZE][SIZE];
  for (i = 0; i < SIZE; i++)
  {
    for (j = 0; j < SIZE; j++)
    {
      x[j][i] = i + j;
    }
  }
}

int main()
{
  //option_one();
  //for (int i = 0; i < 1000; i++)
  //{
  option_two();
  //}
  return 0;
}
