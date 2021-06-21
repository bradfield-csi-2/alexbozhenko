#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <stdbool.h>

int previous_char = '\0';
int current_char = '\0';

int main(void)
{
    bool IS_INSIDE_COMMENT = false;
    bool IS_INSIDE_DOUBLE_QUOTED_STRING = false;

    while ((current_char = getchar()) != EOF)
    {

        if (!IS_INSIDE_COMMENT && current_char == '"')
        {
            IS_INSIDE_DOUBLE_QUOTED_STRING = !IS_INSIDE_DOUBLE_QUOTED_STRING;
        }

        if (!IS_INSIDE_DOUBLE_QUOTED_STRING)
        {
            if (previous_char == '/' && current_char == '*')
            {
                IS_INSIDE_COMMENT = true;
            }
            if (previous_char == '*' && current_char == '/')
            {
                IS_INSIDE_COMMENT = false;
                continue;
            }

            if (!IS_INSIDE_COMMENT && current_char != '/')
            {
                if (previous_char == '/')
                {
                    putchar(previous_char);
                }
                putchar(current_char);
            }
        }
        else
        {
            putchar(current_char);
        }

        previous_char = current_char;
    }
}
