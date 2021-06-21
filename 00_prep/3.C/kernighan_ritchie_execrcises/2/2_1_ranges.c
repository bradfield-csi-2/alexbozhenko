#include <limits.h>
#include <float.h>
#include <stdio.h>

int main(void)
{
    printf("char range:\t[%d;%d]\n", CHAR_MIN, CHAR_MAX);   // 2^8  = 1 byte
    printf("short range:\t[%d;%d]\n", SHRT_MIN, SHRT_MAX);  // 2^16 = 2 bytes
    printf("int range:\t[%d;%d]\n", INT_MIN, INT_MAX);      // 2^32 = 4 bytes
    printf("long range:\t[%ld;%ld]\n", LONG_MIN, LONG_MAX); // 2^64 = 8 bytes

    printf("uchar range:\t[0;%d]\n", UCHAR_MAX);  // 2^8  = 1 byte
    printf("ushort range:\t[0;%d]\n", USHRT_MAX); // 2^16 = 2 bytes
    printf("uint range:\t[0;%u]\n", UINT_MAX);    // 2^32 = 4 bytes
    printf("ulong range:\t[0;%lu]\n", ULONG_MAX); // 2^64 = 8 bytes

    printf("float range:\t[%f;%f]\n", -FLT_MAX, FLT_MAX);  // 2^64 = 8 bytes
    printf("double range:\t[%f;%f]\n", -DBL_MAX, DBL_MAX); // 2^64 = 8 bytes
}
