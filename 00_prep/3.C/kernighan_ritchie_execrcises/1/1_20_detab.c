#include <stdio.h>
#include <assert.h>
#include <stdbool.h>

int main(void)
{
    //assert(false && "My first unit test");
    int chr;
    int TABSTOP = 8;
    int char_pos = 0;
    while ((chr = getchar()) != EOF)
    {
        if (chr == '\t')
        {
            int tab_char_pos = char_pos;
            for (int i = 0; i < (TABSTOP - (tab_char_pos % TABSTOP)); i++)
            {
                putchar(' ');
                char_pos++;
            }
        }
        else
        {
            putchar(chr);
            char_pos++;
        }

        if (chr == '\n')
        {
            char_pos = 0;
        }
    }
}
