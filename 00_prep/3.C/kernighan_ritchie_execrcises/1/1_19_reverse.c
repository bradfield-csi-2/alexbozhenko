#include <stdio.h>
#include <wchar.h>

#define MAXLINE 1000 /* maximum input line size */

void reverse(char line[], int length)
{
    int c;
    int temp;
    if (line[length - 1] == '\n')
    {
        length--;
    }
    for (int i = 0; i < length / 2; i++)
    {
        temp = line[i];
        line[i] = line[length - i - 1];
        line[length - i - 1] = temp;
    }
}

int bgetline(char line[], int lim)
{
    int c;
    int line_pos = 0;

    while ((c = getchar()) != EOF)
    {
        line[line_pos++] = c;
        if (c == '\n')
        {
            line[line_pos] = '\0';
            break;
        }
    }
    return line_pos;
}

void print_line(char line[])
{
    int line_pos = 0;
    int c;
    while ((c = line[line_pos++]) != '\0')
    {
        putchar(c);
    }
}

int main(void)
{
    char line[MAXLINE];
    int current_length;

    while ((current_length = bgetline(line, 0)) != 0)
    {
        reverse(line, current_length);
        print_line(line);
    }
}