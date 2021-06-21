#include <stdio.h>

#define UPPER 300
#define LOWER 0
#define STEP 20

float fahr_to_celsius(float fahr)
{
    return (fahr * (5.0 / 9.0) - 32);
}

void main()
{
    float celsius;

    printf("%3s %6s", "C",
           "F\n");
    for (float fahr = LOWER; fahr <= UPPER; fahr = fahr + STEP)
    {
        celsius = fahr_to_celsius(fahr);
        printf("%3.0f %6.1f\n", fahr, celsius);
    }

    float fahr;

    printf("%6s %3s", "F‭", "☭\n");
    for (float celsius = 134.7; celsius >= -32; celsius = celsius - 11.1)
    {
        fahr = (celsius + 32) * (9.0 / 5.0);
        printf("%6.1f %3f\n", celsius, fahr);
    }
}