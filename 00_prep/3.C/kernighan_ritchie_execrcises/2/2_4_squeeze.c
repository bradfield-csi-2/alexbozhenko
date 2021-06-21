#include <stdio.h>
#include <assert.h>
#include <stdbool.h>
#include <string.h>
#include <ctype.h>

void squeeze(char s1[], char s2[]);

void squeeze(char s1[], char s2[])
{
    size_t read_pos, write_pos;

    for (read_pos = write_pos = 0; read_pos < strlen(s1); read_pos++)
    {
        bool should_squeeze_this_char = false;
        for (size_t s2_pos = 0; s2_pos < strlen(s2); s2_pos++)
        {
            if (s1[read_pos] == s2[s2_pos])
            {
                should_squeeze_this_char = true;
                break;
            }
        }
        if (!should_squeeze_this_char)
        {
            s1[write_pos++] = s1[read_pos];
        }
    }
    s1[write_pos++] = '\0';
}

int main(void)
{
    char s[] = "Cool";
    squeeze(s, "lol");
    assert(strcmp(s, "C") == 0);

    char s2[] = "Test";
    squeeze(s2, "");
    assert(strcmp(s2, "Test") == 0);
    printf("%s", "All tests passed");
}
