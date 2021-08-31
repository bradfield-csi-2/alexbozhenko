#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/wait.h>
#include <unistd.h>
#include "types.h"

int START = 2, END = 20;
char *TESTS[] = {"brute_force", "brutish", "miller_rabin"};
int num_tests = sizeof(TESTS) / sizeof(char *);

// PLAN:
// 1. create queues in suboptimus
// 2. create request message with request for each algorithm, for each
//    number in the range of numbers
// 3. send messages to request queue
// 4. fork and exec N primality.exe workers, where N=number of cpu cores
// 5. receive request messages in primality.exe
// 6. run primality algorithm and publish result to response queue in primality.exe
// 7. receive response messages in suboptimus.exes
// 8. (optionally sort before printing) print results

int main(void)
{
  int testfds[num_tests][2];
  int resultfds[num_tests][2];
  int result, i;
  long n;
  pid_t pid;

  for (i = 0; i < num_tests; i++)
  {
    pipe(testfds[i]);
    pipe(resultfds[i]);

    pid = fork();

    if (pid == -1)
    {
      fprintf(stderr, "Failed to fork\n");
      exit(-1);
    }

    if (pid == 0)
    {
      // we are the child, connect the pipes correctly and exec!
      close(testfds[i][1]);
      close(resultfds[i][0]);
      dup2(testfds[i][0], STDIN_FILENO);
      dup2(resultfds[i][1], STDOUT_FILENO);
      execl("primality.exe", "primality.exe", TESTS[i], (char *)NULL);
    }

    // we are the parent
    close(testfds[i][0]);
    close(resultfds[i][1]);
  }

  // for each number, run each test
  for (n = START; n <= END; n++)
  {
    for (i = 0; i < num_tests; i++)
    {

      // we are the parent, so send test case to child and read results
      write(testfds[i][1], &n, sizeof(n));
      read(resultfds[i][0], &result, sizeof(result));
      printf("%15s says %ld %s prime\n", TESTS[i], n, result ? "is" : "IS NOT");
    }
  }
}
