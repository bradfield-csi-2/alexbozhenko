#include <stdio.h>
#include <assert.h>
#include <stdbool.h>
#include <string.h>
#include <ctype.h>

long my_strpbrk(char s1[], char s2[]);

long my_strpbrk(char s1[], char s2[])
{
    long found_pos = -1;

    for (size_t s1_pos = 0; s1_pos < strlen(s1); s1_pos++)
    {
        for (size_t s2_pos = 0; s2_pos < strlen(s2); s2_pos++)
        {
            if (s1[s1_pos] == s2[s2_pos])
            {
                found_pos = (long)s1_pos;

                return found_pos;
            }
        }
    }
    return found_pos;
}

int main(void)
{
    assert(my_strpbrk("Cool", "lol") == 1);
    assert(my_strpbrk("CHP", "P") == 2);
    assert(my_strpbrk("Cool", "") == -1);
    assert(my_strpbrk("Ventura", "Hi") == -1);

    printf("%s", "All tests passed");
}
