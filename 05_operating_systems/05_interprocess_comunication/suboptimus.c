#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/wait.h>
#include <unistd.h>
#include <sys/sysinfo.h>
#include <mqueue.h>
#include "types.h"

int START = 2, END = 20;
char *TESTS[] = {"brute_force", "brutish", "miller_rabin"};
int num_tests = sizeof(TESTS) / sizeof(char *);

// PLAN:
//+ 1. create queues in suboptimus
//+ 2. create request message with request for each algorithm, for each
//     number in the range of numbers
//+ 3. send messages to request queue
//+ 4. fork and exec N primality.exe workers, where N=number of cpu cores
// 5. receive request messages in primality.exe
// 6. run primality algorithm and publish result to response queue in primality.exe
//+ 7. receive response messages in suboptimus.exes
// 8. (optionally sort before printing)
//+ 9. print results

int main(void)
{
  long n;
  pid_t pid;

  mqd_t request_queue = mq_open(request_queue_name, O_RDWR);
  mqd_t response_queue = mq_open(response_queue_name, O_RDWR);

  for (int i = 0; i < get_nprocs(); i++)
  {
    pid = fork();

    if (pid == -1)
    {
      fprintf(stderr, "Failed to fork\n");
      exit(-1);
    }

    if (pid == 0)
    {
      // we are the child, run primality.exe
      execl("primality.exe", "primality.exe", (char *)NULL);
    }
  }

  // for each number, run each test
  for (n = START; n <= END; n++)
  {
    for (int i = 0; i < NUM_ALGORITHMS; i++)
    {
      // we are the parent, so send test case to child and read results
      struct request req = {.number = n, .alg = i};
      mq_send(request_queue, (const char *)&req,
              sizeof(req), (unsigned int)0);
    }
  }
  for (int i = 0; i <= NUM_ALGORITHMS * (END - START + 1); i++)
  {
    struct response *resp = {0};
    mq_receive(response_queue, (char *)resp,
               sizeof(struct response), (unsigned int)0);
    printf("%15s says %ld %s prime\n",
           ALGORITHMS_STRING[i],
           resp->number,
           resp->result ? "is" : "IS NOT");
  }
}