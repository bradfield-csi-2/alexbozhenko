"Distributed" key-value store built for educational purposes.

Usage:
1. Start 3 server nodes in separate windows, each responsible for 
```
server$ go run . primary_partition 0
server$ go run . primary_partition 1
server$ go run . primary_partition 2
```

2. Start clients in other windows:
```
client$ go run .
Welcome to the distributed K-V store client
We support the following syntax:
get foo
set foo=bar

In [1]: get foo
Out[1]: foofoo
...
```
You can also run expect script:
```
client$ go build; expect expect_script_primary.exp 8000
```
Or do it in parallel with `clush`:
```
client$ go build; clush -f 100  -R exec -w [1-10] expect expect_script_primary.exp 8000
```

In server logs you will see that requests are redirected using consistent hashing algorithm.