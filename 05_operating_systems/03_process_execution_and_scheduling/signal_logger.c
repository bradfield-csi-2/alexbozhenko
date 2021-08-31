#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void signal_handler(int sig)
{
    printf("Caught signal %d:%s\n", sig, strsignal(sig));
    if (sig == SIGINT)
    {
      //  exit(0);
    }
}

int main()
{
    for (int sig = 1; sig <= 31; sig++)
    {
        if (signal(sig, signal_handler) == SIG_ERR)
            printf("signal %d:%s error\n", sig, strsignal(sig));
    }

    while (1)
    {
        /* Wait for the receipt of a signal */
    }
    return 0;
}