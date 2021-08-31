#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/wait.h>
#include <sys/prctl.h>
#include <unistd.h>
#include <sys/sysinfo.h>
#include <mqueue.h>
#include <errno.h>
#include <string.h>
#include "types.h"

int START = 23000000, END = 23000500;
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
  int CPU_CORES = get_nprocs();
  CPU_CORES = 1;

  struct mq_attr attr;

  // TODO do we need this?
  /* initialize the queue attributes */

  // set the value first
  //sudo sysctl  fs.mqueue.msg_max=10240

  //check it:
  //sysctl -a 2>/dev/null | grep mqueue

  // and then run with:
  // prlimit --nofile=20000 ./suboptimus.exe
  attr.mq_flags = 0;
  attr.mq_maxmsg = 1000;
  attr.mq_msgsize = sizeof(struct request);
  attr.mq_curmsgs = 0;

  mqd_t request_queue = mq_open(request_queue_name, O_RDWR | O_CREAT,
                                0644, &attr);
  if (request_queue == (mqd_t)-1)
  {
    printf("%s\n", strerror(errno));
    exit(EXIT_FAILURE);
  }
  attr.mq_msgsize = sizeof(struct response);
  mqd_t response_queue = mq_open(response_queue_name, O_RDWR | O_CREAT,
                                 0644, &attr);

  if (request_queue == (mqd_t)-1)
  {
    printf("%s\n", strerror(errno));
    exit(EXIT_FAILURE);
  }

  for (int i = 0; i < CPU_CORES; i++)
  {
    pid = fork();
    //TODO: undisable fork
    //pid = 1;

    if (pid == -1)
    {
      fprintf(stderr, "Failed to fork\n");
      exit(-1);
    }

    if (pid == 0)
    {
      prctl(PR_SET_PDEATHSIG, SIGKILL);
      // we are the child, run primality.exe
      execl("primality.exe", "primality.exe", (char *)NULL);
    }
  }

  // for each number, run each test
  for (n = START; n <= END; n++)
  {
    for (int i = 0; i < NUM_ALGORITHMS; i++)
    {
      struct request req = {.number = n, .alg = i};

      mq_send(request_queue, (const char *)&req,
              sizeof(req), (unsigned int)0);
      printf("Enqueued %ld with %15s\n",
             req.number,
             ALGORITHMS_STRING[req.alg]);
    }
  }
  struct response resp;
  for (int i = 0; i <= NUM_ALGORITHMS * (END - START + 1); i++)
  {
    mq_receive(response_queue, (char *)&resp,
               sizeof(struct response), (unsigned int)0);
    printf("%15s says %ld %s prime. Took %2.10f s.\n",
           ALGORITHMS_STRING[resp.alg],
           resp.number,
           resp.result ? "is" : "IS NOT",
           resp.duration);
  }
}