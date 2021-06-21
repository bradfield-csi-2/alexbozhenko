#include "vendor/unity.h"
#include <stdlib.h>

extern int fib(int n);
extern int index(int *matrix, int rows, int cols, int rindex, int cindex);
extern void transpose(int *in, int *out, int rows, int cols);
extern float volume(float radius, float height);

void setUp(void)
{
}

void tearDown(void)
{
}

void test_fib_0(void) { TEST_ASSERT_EQUAL(0, fib(0)); }
void test_fib_1(void) { TEST_ASSERT_EQUAL(1, fib(1)); }
void test_fib_2(void) { TEST_ASSERT_EQUAL(1, fib(2)); }
void test_fib_3(void) { TEST_ASSERT_EQUAL(2, fib(3)); }
void test_fib_4(void) { TEST_ASSERT_EQUAL(3, fib(4)); }
void test_fib_5(void) { TEST_ASSERT_EQUAL(5, fib(5)); }
void test_fib_10(void) { TEST_ASSERT_EQUAL(55, fib(10)); }
void test_fib_12(void) { TEST_ASSERT_EQUAL(144, fib(12)); }

void test_index_row(void)
{
  int matrix[1][4] = {{1, 2, 3, 4}};
  TEST_ASSERT_EQUAL(3, index((int *)matrix, 1, 4, 0, 2));
}

void test_index_col(void)
{
  int matrix[4][1] = {{1}, {2}, {3}, {4}};
  TEST_ASSERT_EQUAL(2, index((int *)matrix, 4, 1, 1, 0));
}

void test_index_rect(void)
{
  int matrix[2][3] = {{1, 2, 3}, {4, 5, 6}};
  TEST_ASSERT_EQUAL(6, index((int *)matrix, 2, 3, 1, 2));
}

void test_cone_volume_0_0(void)
{
  TEST_ASSERT_FLOAT_WITHIN(0.01, 0.0, volume(0.0, 0.0));
}
void test_cone_volume_1_2(void)
{
  TEST_ASSERT_FLOAT_WITHIN(0.01, 2.09, volume(1.0, 2.0));
}
void test_cone_volume_55_55(void)
{
  TEST_ASSERT_FLOAT_WITHIN(0.01, 174.23, volume(5.5, 5.5));
}
void test_cone_volume_1234_5678(void)
{
  TEST_ASSERT_FLOAT_WITHIN(0.01, 9.05, volume(1.234, 5.678));
}

int main(void)
{
  UNITY_BEGIN();

  RUN_TEST(test_fib_0);
  RUN_TEST(test_fib_1);
  RUN_TEST(test_fib_2);
  RUN_TEST(test_fib_3);
  RUN_TEST(test_fib_4);
  RUN_TEST(test_fib_5);
  RUN_TEST(test_fib_10);
  RUN_TEST(test_fib_12);

  RUN_TEST(test_index_row);
  RUN_TEST(test_index_col);
  RUN_TEST(test_index_rect);

  RUN_TEST(test_cone_volume_0_0);
  RUN_TEST(test_cone_volume_1_2);
  RUN_TEST(test_cone_volume_55_55);
  RUN_TEST(test_cone_volume_1234_5678);

  return UNITY_END();
}
