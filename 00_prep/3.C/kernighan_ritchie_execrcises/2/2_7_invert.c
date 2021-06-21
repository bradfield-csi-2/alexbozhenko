#include <stdio.h>
#include <assert.h>
#include <stdbool.h>
#include <string.h>
#include <ctype.h>
#include "lib/bin_printer.h"

int rightrot(int x, int p, int n);

int rightrot(int x, int p, int n)
{

    unsigned bitmask_with_n_bits_from_p_turned_on = ~((~(~(unsigned)0 << (p + 1))) ^ (~(unsigned)0 << (p + 1 - n)));

    printf("%-25s %40s\n", "bitmask=", to_bin_string((int)bitmask_with_n_bits_from_p_turned_on));
    printf("%-25s %40s\n", "x=", to_bin_string(x));
    int x_with_subset_bits_inverted = x ^ (int)bitmask_with_n_bits_from_p_turned_on;

    return x_with_subset_bits_inverted;
}

int main(void)
{
    test_to_bin_string();

    int x = 55880066;
    int p = 9;
    int n = 6;

    printf("%-25s %40s\n", "result=", to_bin_string(rightrot(x, p, n)));
}
