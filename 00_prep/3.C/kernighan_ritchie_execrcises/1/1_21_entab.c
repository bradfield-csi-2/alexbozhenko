#include <stdio.h>
#include <assert.h>
#include <stdbool.h>

#define MAXLINE 40 /* maximum input line size */
#define TABSIZE 8

int my_getline(char line[], int limit);
bool strings_equal(char s1[], char s2[]);
void entab_line(char line[], char entabbed[]);
void test_entab_line(char orig[], char entabbed[], char description[]);

/* read up to limit characters into line, return length */
int my_getline(char line[], int limit)
{
    int chr;
    int i;
    for (i = 0; i < limit && (chr = getchar()) != EOF && chr != '\n'; i++)
    {
        line[i] = (char)chr;
    }
    if (line[i] == '\n')
    {
        line[i++] = '\n';
    }
    line[i] = '\0';
    return i;
}

bool strings_equal(char s1[], char s2[])
{
    bool equal = true;
    for (int i = 0; i < MAXLINE; i++)
    {
        if (s1[i] != s2[i])
        {
            equal = false;
        }
        if (s1[i] == '\0' ||
            s2[i] == '\0')
        {
            break;
        }
    }
    return equal;
}

void entab_line(char line[], char entabbed[])
{
    int line_pos;
    int entabbed_pos = 0;
    int consequtive_space_seen = 0;
    for (line_pos = 0; line[line_pos] != '\0'; line_pos++)
    {
        if (line[line_pos] == ' ')
        {
            consequtive_space_seen++;

            if ((line_pos % TABSIZE == (TABSIZE - 1)))
            {
                entabbed[entabbed_pos++] = '\t';
                consequtive_space_seen = 0;
            }
        }
        else
        {
            for (int i = 0; i < consequtive_space_seen; i++)
            {
                entabbed[entabbed_pos++] = ' ';
            }
            consequtive_space_seen = 0;
            entabbed[entabbed_pos++] = line[line_pos];
        }
    }

    // CR abozhenko: how to dedup?
    for (int i = 0; i < consequtive_space_seen; i++)
    {
        entabbed[entabbed_pos++] = ' ';
    }
    entabbed[entabbed_pos++] = '\0';
}

void test_entab_line(char orig[], char entabbed[], char description[])
{
    char test_entabbed_line[MAXLINE];
    entab_line(orig, test_entabbed_line);
    assert(strings_equal(entabbed, test_entabbed_line) && description);
}

int main(void)
{
    assert(!strings_equal("abcd\n", "defg\n") && "not equal");
    assert(!strings_equal("abcd", "abcde") && "different length");
    assert(strings_equal("abcd", "abcd") && "equal");
    assert(!strings_equal("", "abcd") && "empty string");

    test_entab_line("a       b", "a\tb", "one tab between chars");
    test_entab_line("", "", "empty string");
    test_entab_line("        a", "\ta", "single tab at the beginning");
    test_entab_line("123456  ", "123456\t", "two spaces before tab");
    test_entab_line("1234567 ", "1234567\t", "one space before tab");
    test_entab_line("12345678", "12345678", "no tab");
    test_entab_line("12345678 ", "12345678 ", "no tab");
    test_entab_line("12345678  ", "12345678  ", "no tab");
    test_entab_line("12345678   ", "12345678   ", "no tab");

    char line[MAXLINE];
    char entabbed[MAXLINE];

    while (my_getline(line, MAXLINE) != 0)
    {
        entab_line(line, entabbed);
        printf("%s\n", entabbed);
    }
}
