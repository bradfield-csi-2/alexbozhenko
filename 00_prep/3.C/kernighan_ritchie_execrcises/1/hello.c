#include <stdio.h>
#include <zlib.h>
#include <string.h>

enum escapes
{
    BELL = '\a',
    BACKSPACE = '\b',
    TAB = '\t',
    NEWLINE = '\n',
    VTAB = '\v',
    RR = '\r'
};

enum months
{
    JAN = 1,
    FEB = 1,
    MAR,
    APR,
    MAY,
    JUN,
    JUL,
    AUG,
    SEP,
    OCT,
    NOV,
    DEC
};

void main()
{
    printf("hello, ", JAN,
           " world!\t!\n\r \v tt \x30 ");
}
