#define _DEFAULT_SOURCE
#include <dirent.h>
#include <stdio.h>
#include <stdlib.h>
#include <ctype.h>
#include <stdio.h>
#include <math.h>
#include <stdbool.h>

#define SIZE 10

int getfloat(float *pn, FILE *stream);

int getfloat(float *pn, FILE *stream)
{
    int c;
    float sign;
    while (isspace(c = getc(stream)))
        ;

    if (!isdigit(c) && c != EOF && c != '+' && c != '-')
    {
        // The following line is present in the book, but
        // imho it should not. With that line, c that is not digit
        // will be pushed back and force to stream
        // on every call to this function
        //ungetc(c, stream);
        return 0; //Not a float
    }
    sign = (c == '-') ? -1.0 : 1.0;
    if (c == '+' || c == '-')
    {
        c = getc(stream);
    }
    bool found_float = false;
    for (*pn = 0.0; isdigit(c); c = getc(stream))

    {
        *pn = *pn * 10 + (float)(c - '0');
        found_float = true;
    }
    if (c == '.')
    {
        c = getc(stream);
    }

    float fractional_part;
    long double exp;
    for (fractional_part = 0, exp = 1; isdigit(c); c = getc(stream), exp++)
    {
        fractional_part = fractional_part + ((float)(c - '0') / (float)powl(10, exp));
        found_float = true;
    }
    *pn = *pn + fractional_part;
    *pn = *pn * sign;

    if (c != EOF)
    {
        ungetc(c, stream);
        return c;
    }
    else
    {
        if (found_float == true)
        {
            //ok, when we parsed a valid number in this call, but also reached
            // EOF, we want to return something rather than EOF to make sure
            // array size is correctly incremented by the caller.
            return 22;
        }
        else
        {
            return c;
        }
    }
}

int main(void)
{
    int n;
    float array[SIZE] = {0.0};

    FILE *test_stream;
    test_stream = fopen("/tmp/test_stream", "w+");
    fputs("88 -98.6 d d 42.34", test_stream);
    fclose(test_stream);

    test_stream = fopen("/tmp/test_stream", "r");

    for (n = 0; n < SIZE && getfloat(&array[n], test_stream) != EOF; n++)
        ;

    fclose(test_stream);

    for (int i = 0; i < n; i++)
    {
        printf("array[%d]=%f\n", i, array[i]);
    }
}
