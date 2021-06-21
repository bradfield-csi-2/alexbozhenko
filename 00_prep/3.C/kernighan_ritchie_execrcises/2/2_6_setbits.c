#include <stdio.h>
#include <assert.h>
#include <stdbool.h>
#include <string.h>
#include <ctype.h>
#include "lib/bin_printer.h"

int setbits(int x, int p, int n, int y);

int setbits(int x, int p, int n, int y)
{

    unsigned bitmask_with_n_bits_from_p_turned_off = (~(~(unsigned)0 << (p + 1))) ^ (~(unsigned)0 << (p + 1 - n));

    printf("%-25s %40s\n", "bitmask=", to_bin_string((int)bitmask_with_n_bits_from_p_turned_off));
    printf("%-25s %40s\n", "x=", to_bin_string(x));
    int x_without_bits_subset = x & (int)bitmask_with_n_bits_from_p_turned_off;
    printf("%-25s %40s\n", "x_without_bits_subset=", to_bin_string(x_without_bits_subset));

    printf("%-25s %40s\n", "y=", to_bin_string(y));
    int y_bits_subset_only = (y << (p + 1 - n)) & ~(int)bitmask_with_n_bits_from_p_turned_off;
    printf("%-25s %40s\n", "y_bits_subset_only", to_bin_string(y_bits_subset_only));

    int result = x_without_bits_subset ^ y_bits_subset_only;
    return result;
}

int main(void)
{
    test_to_bin_string();

    int x = 82397679;
    int p = 9;
    int n = 6;
    int y = 68979;

    printf("%-25s %40s\n", "result=", to_bin_string(setbits(x, p, n, y)));
}
