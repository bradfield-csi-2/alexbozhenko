Poking around with linux IPCs:
Implementing primality testing with two processes communicating
with each other using [POSIX message queues](https://man7.org/linux/man-pages/man7/mq_overview.7.html)

`primality.exe` creates two queues: requests and responses, and then forks NUM_CPUs processes, and there execs primality.exe workers.
Then it forks itself to do the enqueuing in a separate process, to make sure we do not block the main process, that is trying to consume responses.
Eeach worker is trying to consume from requests queue, performing a check using specified algorithm, and publishing response to the responses queue.
Finally, `primality.exe` consumes all the results and prints them to `stdout`.

To compile:
```
gcc  -lm -lrt  primality.c -o primality.exe
gcc -lrt  suboptimus.c -o suboptimus.exe
```

Cleaning up of the queue is not implemented, so I was manually cleaning up the queues like this:
```
rm -rf /dev/mqueue/*;  ./suboptimus.exe 
```