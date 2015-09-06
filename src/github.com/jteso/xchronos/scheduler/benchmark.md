PriorityQueue
```
BenchmarkPQueue5000          500       2339640 ns/op
BenchmarkPQueue50000          50      23719483 ns/op
BenchmarkPQueue500000          5     291451387 ns/op
BenchmarkPQueue1000000         2     667503969 ns/op
```

CalendarQueue (width=360, len=2)
```
BenchmarkCalq5000           1000       2315311 ns/op
BenchmarkCalq50000            50      24344656 ns/op
BenchmarkCalq500000            5     265413434 ns/op
BenchmarkCalq1000000           2     578307656 ns/op
```


CalendarQueue (width=36000, len=100)
```
BenchmarkCalq5000        500       2307071 ns/op
BenchmarkCalq50000        50      23949934 ns/op
BenchmarkCalq500000        5     288294056 ns/op
BenchmarkCalq100000       20     664841857 ns/op
```

CalendarQueue (width=36000, len=1000)
```
BenchmarkCalq5000            500       2754834 ns/op
BenchmarkCalq50000            50      25626809 ns/op
BenchmarkCalq500000            5     303771504 ns/op
BenchmarkCalq1000000           2     629922164 ns/op
```