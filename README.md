## Overview
check tests
```shell
make test
```
up in docker 
```shell
make run
```
target url us ws://localhost:8000/ws/duty_processor

### Implementation
* main package is `internal/service/duty`. package separate into 2 parts - orchestrator and executor
  * `orchestrator` - responsible for route incoming messages to specific executor
  * `executor` - responsible for process incoming messages for validator. so it allow process messages in parallel per validator
* `internal/service/duty/executor` package store incoming messages into map, where key is `DutyType` and value is list of duty requests.
  * background goroutine process this map. so it allows to run process in parallel per validator for each duty type
  * list allow to keep message order by store it in order of height
  * list is used as execution stack - only one execution at time

### Improvements
1. current server does not support multiply handlers. if 2 or more clients will connect to server and will try send duty requests per validator - some improvements need to be done. Orchestrator should re-stream responses to each client and track requests which each client sent.
2. deduplication does not implemented. if client will send same duty request twice - it will be processed twice. some external storage should be used for keep track of processed requests