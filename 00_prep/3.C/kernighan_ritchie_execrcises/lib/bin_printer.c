#include <stdio.h>
#include <assert.h>
#include <stdbool.h>
#include <string.h>
#include <ctype.h>

#define INT_BITS 32
#define BYTE_SIZE 8
#define BYTE_SPACERS 3

char *to_bin_string(int n);
void test_to_bin_string(void);

char *to_bin_string(int n)
{
    static char bits_string[INT_BITS + 1];
    int bit_pos_inside_string;

    for (bit_pos_inside_string = INT_BITS - 1; bit_pos_inside_string >= 0; bit_pos_inside_string--)
    {
        bits_string[bit_pos_inside_string] = (n & 1) ? '1' : '0';
        n = n >> 1;
    }
    bits_string[INT_BITS] = '\0';
    static char formatted_bits[INT_BITS + 1 + BYTE_SPACERS];
    sprintf(formatted_bits, "%.*s %.*s %.*s %.*s",
            BYTE_SIZE, bits_string,
            BYTE_SIZE, bits_string + BYTE_SIZE,
            BYTE_SIZE, bits_string + 2 * BYTE_SIZE,
            BYTE_SIZE, bits_string + 3 * BYTE_SIZE);
    return (formatted_bits);
}

void test_to_bin_string(void)
{
    // printf("-2\t\t %s\n", to_bin_string(-2));
    // printf("-1\t\t %s\n", to_bin_string(-1));
    // printf("0\t\t %s\n", to_bin_string(0));
    // printf("1\t\t %s\n", to_bin_string(1));
    // printf("2\t\t %s\n", to_bin_string(2));
    // printf("~0\t\t %s\n", to_bin_string(~0));
    // printf("2147483647\t %s\n", to_bin_string(2147483647));
    // printf("-2147483647\t %s\n", to_bin_string(-2147483647));

    assert(strcmp("11111111 11111111 11111111 11111110", to_bin_string(-2)) == 0);
    assert(strcmp("11111111 11111111 11111111 11111111", to_bin_string(-1)) == 0);
    assert(strcmp("00000000 00000000 00000000 00000000", to_bin_string(0)) == 0);
    assert(strcmp("00000000 00000000 00000000 00000001", to_bin_string(1)) == 0);
    assert(strcmp("00000000 00000000 00000000 00000010", to_bin_string(2)) == 0);
    assert(strcmp("11111111 11111111 11111111 11111111", to_bin_string(~0)) == 0);
    assert(strcmp("01111111 11111111 11111111 11111111", to_bin_string(2147483647)) == 0);
    assert(strcmp("10000000 00000000 00000000 00000001", to_bin_string(-2147483647)) == 0);
}
