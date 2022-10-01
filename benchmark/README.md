# High-level plan for benchmarking

- [X] Initial benchmark setup
- [X] prepare random rule-engine-config generator and valid input generator
- [ ] prepare benchmarking scripts

## Command

```bash
    cd ./benchmark/
    go test -bench=. -benchtime=5s -benchmem
```

## Result History

1. [commit-31b10ac](https://github.com/niharrathod/ruleengine-core/tree/31b10acc10f4bb157b21ca66d93894234f1738f4)
```
goos: linux
goarch: amd64
pkg: github.com/niharrathod/ruleengine-core/benchmark
cpu: AMD Ryzen 5 5600X 6-Core Processor             
Benchmark_RuleEngine_10_Rule_EvaluateComplete/validInput-12               582061              9339 ns/op            4068 B/op         60 allocs/op
Benchmark_RuleEngine_10_Rule_EvaluateComplete/invalidInput-12             623803              9194 ns/op            4067 B/op         60 allocs/op
Benchmark_RuleEngine_100_Rule_EvaluateComplete/validInput-12               96234             63249 ns/op           19277 B/op        420 allocs/op
Benchmark_RuleEngine_100_Rule_EvaluateComplete/invalidInput-12             96860             61492 ns/op           19275 B/op        420 allocs/op
Benchmark_RuleEngine_1000_Rule_EvaluateComplete/validInput-12               8499            621958 ns/op          171462 B/op       4024 allocs/op
Benchmark_RuleEngine_1000_Rule_EvaluateComplete/invalidInput-12             8935            579088 ns/op          171377 B/op       4020 allocs/op
PASS
ok      github.com/niharrathod/ruleengine-core/benchmark        35.450s
```
