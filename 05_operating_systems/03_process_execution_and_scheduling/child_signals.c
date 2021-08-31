#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/wait.h>

void main() {
    printf("start\n");
    int rc = fork();
    if (rc < 0) {
        printf("epic fail\n");
        exit(1);
    }
    if (rc == 0) {
        printf("im child. pid = %d\n", getpid());
        execlp("sleep", "sleep", "5", (char *) NULL);
    }
    else {
        int* wstatus = malloc(sizeof *wstatus);
        int rc_wait = waitpid(rc, wstatus, 0);
        if WIFEXITED (*wstatus) {
            printf("im parent. pid = %d, child terminated normally. It's pid was %d\n", getpid(), rc);
        } else if WIFSIGNALED(*wstatus) {
            printf("im parent. pid = %d, child terminated due to a signal. It's pid was %d\n", getpid(), rc);
            printf("Signal that terminated was %d\n", WTERMSIG(*wstatus));
        }
        free(wstatus);
    }

}