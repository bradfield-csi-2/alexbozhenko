# Objective:
Design a system where user-supplied URLs are turned into “shortened” versions of the original.

## Functional requirements
* 10 million new links created per day
* Links will not expire
* To create url, users need to sign up 
* 1 generated link will result in 100 reads on average over lifetime of the URL
* Capture as much data as possible about users activity

## Non-functional requirements
* Low redirect latency
* Extremely high availability target


## Estimations
### Write RPS:
<!-- http://www.sciweavers.org/free-online-latex-equation-editor -->
```
      7                       7                             
1 * 10  URLs           1  * 10  URLs          2             
-------------    ~    --------------- = 1 * 10   =  100 write RPS
86400 seconds                5                               
                       1 * 10  seconds                       
```

### How many links will be created over 1 year?
```
      7                   7           2           9
1 * 10   *  365  ~  1 * 10   *  4 * 10   =  4 * 10   links = 4 billion links
```

### How many reads we will have over a year?
Assuming 80% of link reads will happen in the first year,
without taking into account reads from "old" urls  :
```
      9                     10           11                      
4 * 10   *  80  =  32  *  10   ~ 3  *  10    =  300 billion reads
```
or
```
       11              11                        
 3 * 10          3 * 10            4             
----------- ~ --------------  =  10   =  1000 read RPS
365 * 86400         2      5                     
              3 * 10   * 10                      
```

## Total storage for links over one year:
Assuming short link + URL, on average, will take 1KB:
```
      9                  9        3                  12                       
4 * 10    *  1KB ~ 4 * 10    *  10  Bytes   =  4 * 10   Bytes  =  4 Terabytes
```

## Storage for logging activity over one year:
Assuming that we log 4KB per read request:
```
      11                 11             3                    14                15
3 * 10    *  4KB ~ 3 * 10    *  4  *  10  Bytes  =  12  *  10   bytes ~ 1  * 10   Bytes  =  1 Petabyte
```
So if we log everything, and store raw(not aggregated) data, we need at least:
```1PB/4TB = 256 HDDs a year```
to store the read log data.

### Conclusion from estimations
Writes(url+shortened url) may fit on a single server.  
Reads ~ read log data needs to be distributed.
Generated log size won't fit on a single server, and also write IOPS rate would not fit(if we use HDDs, not SSDs).

## Initial design

* Single PostgreSQL master, hosting a db containing mapping of short urls to long urls. Indexed on short urls column.
  * Same server hosts backend app serving creation of the URLs and stats dashboard.
  * `shortener.shadycorp.com` is behind a loadbalancer, pointing to this single server.
  * In case of a crash, any of the replicas can be assigned to be a new master.
* 16 read-only postgresql replicas, each containing entire copy of the db.
  * Each of the replicas run proxy app that, upon receiving a request, looks up the long url in it's local RO postgres replica.
  * This proxy app also logs all the details about the client in local log.
  * Let's say raw logs are truncated(or archived to cold storage) after one year, so such server can have 16 4TB HDDs to fit all the logs.
  * All replicas are behind round-robin load balancer.
* There is map-reduce job that pulls logs from all the proxies, let's say, every 5 minutes, and writes aggregated stats about each url into master PostgreSQL instance. 

