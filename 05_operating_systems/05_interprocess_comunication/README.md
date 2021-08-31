Poking around with linux IPCs:
Implementing primality testing with two processes communicating
with each other using [POSIX message queues](https://man7.org/linux/man-pages/man7/mq_overview.7.html)

To compile:
```
gcc  -lm -lrt  primality.c -o primality.exe
gcc -lrt  suboptimus.c -o suboptimus.exe
```

Cleaning up of the queue is not implemented, so I was manually cleaning up the queues like this:
```
rm -rf /dev/mqueue/*;  ./suboptimus.exe 
```