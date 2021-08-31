#include <stdbool.h>

// from https://stackoverflow.com/a/10966395
#define FOREACH_ALGO(ALGO) \
    ALGO(brute_force)      \
    ALGO(brutish)          \
    ALGO(miller_rabin)

#define GENERATE_ENUM(ENUM) ENUM,
#define GENERATE_STRING(STRING) #STRING,

typedef enum
{
    FOREACH_ALGO(GENERATE_ENUM)
} ALGORITHMS_ENUM;

static const char *ALGORITHMS_STRING[] = {
    FOREACH_ALGO(GENERATE_STRING)};

const int NUM_ALGORITHMS = sizeof ALGORITHMS_STRING / sizeof *ALGORITHMS_STRING;

const char *request_queue_name = "/suboptimus_request_queue";
const char *response_queue_name = "/suboptimus_response_queue";

// typedef enum
// {
//     brute_force,
//     brutish,
//     miller_rabin,
//     // Please, add new algorithm before this comment
//     num_algorithms
// } algorithm;

struct request
{
    long number;
    ALGORITHMS_ENUM alg;
};

struct response
{
    long number;
    ALGORITHMS_ENUM alg;
    double duration;
    bool result;
};