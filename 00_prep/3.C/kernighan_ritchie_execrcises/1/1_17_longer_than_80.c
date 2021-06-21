#include <stdio.h>
#include <wchar.h>

#define MAXLINE 1000 /* maximum input line size */
#define LENGTH_THRESHOLD 80

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

void main()
{
    char line[MAXLINE];
    int current_length;

    while ((current_length = bgetline(line, 0)) != 0)
    {
        if (current_length > LENGTH_THRESHOLD)
        {
            print_line(line);
        }
    }
}