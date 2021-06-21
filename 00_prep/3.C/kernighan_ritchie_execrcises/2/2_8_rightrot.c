#include <stdio.h>
#include <assert.h>
#include <stdbool.h>
#include <string.h>
#include <ctype.h>

#include <limits.h>
#include "lib/bin_printer.h"

unsigned rightrot(unsigned x, int n);

unsigned rightrot(unsigned x, int n)
{
    int bits_in_int = sizeof(int) * CHAR_BIT;
    printf("%d\n", bits_in_int);
    n = n % bits_in_int;
    printf("%d\n", n);

    printf("%-25s %40s\n", "x=", to_bin_string((int)x));
    //unsigned x_rotated = ((x >> n)); // & ~(~(unsigned)0 << n)); // ^ (x << (bits_in_int - 1 - n));
    //unsigned x_rotated = (x << (bits_in_int - n));
    unsigned left_bits_moved_to_end = (x >> n);
    printf("%-25s %40s\n", "left_bits_moved_to_end=", to_bin_string((int)left_bits_moved_to_end));

    unsigned right_bits_moved_to_beginning = (x << (bits_in_int - n));
    printf("%-25s %40s\n", "right_bits_moved_to_beg=", to_bin_string((int)right_bits_moved_to_beginning));
    unsigned x_rotated = left_bits_moved_to_end ^ right_bits_moved_to_beginning;

    return x_rotated;
}

int main(void)
{
    test_to_bin_string();

    unsigned x = 0x712345;
    int n = 20;

    printf("%-25s %40s\n", "result=", to_bin_string((int)rightrot(x, n)));
}
