Educational key-value store.

Let's start with map serialized to disk using Gob, 
and available for `get`s and `set`s via http endpoint

Usage:
1. Start a server in one window:
```
server$ go run .
Welcome to the distributed K-V store serv
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
client$ go build; expect expect_script.exp
```