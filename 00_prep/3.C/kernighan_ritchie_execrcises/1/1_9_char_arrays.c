#include <stdio.h>
#include <wchar.h>

#define MAXLINE 1000 /* maximum input line size */

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

void copy_line(char from[], char to[])
{
    for (int i = 0; from[i] != '\0'; i++)
    {
        to[i] = from[i];
    }
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

void main()
{
    char line[MAXLINE];
    char longest_line[MAXLINE];
    int max_length = 0;
    int current_length = 0;

    while ((current_length = bgetline(line, 0)) != 0)
    {
        if (current_length > max_length)
        {
            max_length = current_length;
            copy_line(line, longest_line);
        }
    }
    print_line(longest_line);
    printf("%d", max_length);
}