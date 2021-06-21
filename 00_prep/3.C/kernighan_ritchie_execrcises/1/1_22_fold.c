#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <stdbool.h>

#define MAXLINE 10000 /* maximum input line size */

void fold_line(char line[], char folded_line[], int max_line_length);

bool is_blank(char c);
bool is_blank(char c)
{
    return (c == ' ' || c == '\t');
}

void fold_line(char line[], char folded_line[], int max_line_length)
{

    int next_newline_insert_position = max_line_length - 1;
    int last_seen_word_end = -1;
    int folded_line_pos = 0;
    bool is_inside_word = !is_blank(line[0]);
    for (int line_pos = 0; line[line_pos] != '\0'; line_pos++)
    {
        if (is_blank(line[line_pos]) && is_inside_word)
        {
            last_seen_word_end = line_pos - 1;
            is_inside_word = false;
        }
        else
        {
            is_inside_word = true;
        }

        folded_line[folded_line_pos++] = line[line_pos];
        if (line_pos == next_newline_insert_position)
        {
            if (last_seen_word_end != -1)
            {
                int last_word_end_offset = line_pos - last_seen_word_end;
                folded_line_pos = folded_line_pos - last_word_end_offset;
                folded_line[folded_line_pos++] = '\n';
                for (int i = last_seen_word_end + 1; i <= line_pos; i++)
                {
                    folded_line[folded_line_pos++] = line[i];
                }
                next_newline_insert_position = last_seen_word_end + max_line_length;
            }
            else
            {
                folded_line[folded_line_pos++] = '\n';
                next_newline_insert_position += max_line_length;
            }
            is_inside_word = false;
            last_seen_word_end = -1;
        }
    }
    folded_line[folded_line_pos++] = '\0';
}

int main(void)
{
    char folded_line[MAXLINE];

    fold_line("0123456789ABCDEF", folded_line, 15);

    assert(
        (strcmp(
             folded_line, "0123456789ABCDE\nF") == 0) &&
        "should insert newline on max_line_length with no blank chars");

    fold_line("0123456789ABCDEF", folded_line, 5);
    assert(
        (strcmp(
             folded_line, "01234\n56789\nABCDE\nF") == 0) &&
        "should insert multiple newlines on very long line with no blank chars");

    fold_line("012345 67 89 ABCDEF", folded_line, 5);
    assert(
        (strcmp(
             folded_line, "01234\n5 67\n 89\n ABCD\nEF") == 0) &&
        "should insert multiple newlines on very long line with blank chars");

    char tmp[MAXLINE] = "___________________________";
    strcpy(folded_line, tmp);

    fold_line("But I must explain to you how all", folded_line, 5);
    assert(
        (strcmp(
             folded_line, "But\n I mu\nst\n expl\nain\n to\n you\n how\n all") == 0) &&
        "should insert multiple newlines on very long line with blank chars");

    char line[MAXLINE];
    scanf("%9999[^\n]", line);
    fold_line(line, folded_line, 5);
    printf("%s\n", folded_line);
}
