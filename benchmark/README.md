# High-level plan for benchmarking

- [X] Initial benchmark setup
- [X] prepare random rule-engine-config generator and valid input generator
- [ ] prepare benchmarking scripts

## Command

```bash
    cd ./benchmark/
    go test -bench=. -benchtime=5s -benchmem
```
