#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>

void main(int argc, char *argv[])
{
    unsigned long long int mem_size_bytes = (unsigned long long int)atoi(argv[1]);
    mem_size_bytes = mem_size_bytes * (unsigned long long int)1024 * (unsigned long long int)1024;
    printf("%s\n", argv[1]);
    printf("%llu\n", mem_size_bytes);
    char *array = malloc(mem_size_bytes * sizeof(char));
    sleep(10);
    for (unsigned long long int i = 0; i < mem_size_bytes; i++)
    {
        if (i % 1000000 == 0)
        {
            printf("%luu\n", i);
        }
        array[i] = i % 128;
        if (i == mem_size_bytes - 1)
        {
            i = 0;
        }
    }
}
