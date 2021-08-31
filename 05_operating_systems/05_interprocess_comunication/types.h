#include <stdbool.h>

typedef enum
{
    brute_force,
    brutish,
    miller_rabin
} algorithm;

struct request
{
    long number;
    algorithm alg;
};

struct response
{
    long number;
    algorithm alg;
    double duration;
    bool result;
};