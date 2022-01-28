Estimate cardinality of a set using 2^n and
[LogLog](https://engineering.fb.com/2018/12/13/data-infrastructure/hyperloglog/) approach.

```
go build && find /usr/share/dict/ -type f | xargs -I% ./maxFNVhashZeroes %
Actual cardianlity         = 127466 /usr/share/dict/british-english
2^n cardinality estimate   = 131072
LogLog cardinality estimate= 129508
=========================
Actual cardianlity         = 88328 /usr/share/dict/finnish
2^n cardinality estimate   = 2097152
LogLog cardinality estimate= 88618
=========================
Actual cardianlity         = 56329 /usr/share/dict/spanish
2^n cardinality estimate   = 65536
LogLog cardinality estimate= 56230
=========================
Actual cardianlity         = 76258 /usr/share/dict/ogerman
2^n cardinality estimate   = 65536
LogLog cardinality estimate= 78135
=========================
Actual cardianlity         = 123115 /usr/share/dict/american-english
2^n cardinality estimate   = 131072
LogLog cardinality estimate= 126219
=========================
Actual cardianlity         = 304736 /usr/share/dict/ngerman
2^n cardinality estimate   = 131072
LogLog cardinality estimate= 319706
=========================
Actual cardianlity         = 189057 /usr/share/dict/catala
2^n cardinality estimate   = 131072
LogLog cardinality estimate= 193016
=========================
Actual cardianlity         = 221377 /usr/share/dict/french
2^n cardinality estimate   = 262144
LogLog cardinality estimate= 221822
=========================
Actual cardianlity         = 92034 /usr/share/dict/italian
2^n cardinality estimate   = 131072
LogLog cardinality estimate= 92479
=========================
```