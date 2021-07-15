Bloom filter implementation

`pprof` was used to remove all unnecessary work:
```
go build && ./bloom_filter   -cpuprofile /tmp/cpuprofile && ll /tmp/cpuprofile && pprof -http 0.0.0.0:9876 /tmp/cpuprofile
```
Graph can be found in `pprof_profile.svg`. 
Runtime of the program is dominated by calculating the has function itself.
