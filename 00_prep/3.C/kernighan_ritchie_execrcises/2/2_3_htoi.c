#include <stdio.h>
#include <assert.h>
#include <stdbool.h>
#include <string.h>
#include <ctype.h>

int htoi(char s[]);
int htoi(char s[])
{
    int result = 0;
    size_t s_length = strlen(s);
    size_t s_pos;
    if ((s_length >= 3) && (s[0] == '0' && (s[1] == 'x' || s[1] == 'X')))
    {
        s_pos = 2;
    }
    else
    {
        s_pos = 0;
    };

    for (; s_pos < s_length; s_pos++)
    {
        int c = tolower(s[s_pos]);
        if (isdigit(c))
        {
            result = (result * 16) + (c - '0');
        }
        else if (isxdigit(c))
        {
            result = (result * 16) + (10 + c - 'a');
        }
    }
    return result;
}

int main(void)
{
    assert((htoi("0x") == 0) && "empty");
    assert((htoi("") == 0) && "empty");
    assert((htoi("0x1F") == 31) && "prefix with x");
    assert((htoi("0X1f") == 31) && "prefix with X");
    assert((htoi("0X0") == 0) && "prefix with X");
    assert((htoi("0x0") == 0) && "prefix with X");
    assert((htoi("1f") == 31) && "without prefix");
    assert((htoi("ff") == 255) && "without prefix");
    assert((htoi("DeAdBEE") == 233495534) && "without prefix");
    // This will overflow and will be detected by -fsanitize=undefined
    // assert((htoi("DeAdBEEF") == 3735928559) && "without prefix");
    printf("%s\n", "All tests passed!");
}
