[![node-graph-flow](https://github.com/WolvenSpirit/node-graph-flow/actions/workflows/go.yml/badge.svg)](https://github.com/WolvenSpirit/node-graph-flow/actions/workflows/go.yml) [![codecov](https://codecov.io/gh/WolvenSpirit/node-graph-flow/branch/main/graph/badge.svg?token=hDXMUdD4L1)](https://codecov.io/gh/WolvenSpirit/node-graph-flow)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/WolvenSpirit/node-graph-flow)

## node graph flow

Simple frame for node graph inspired task processing.

### Disclaimer: 

- Each individual flow can be executed concurrently on separate goroutines.

- At node level, all nodes of an individual graph should execute on the same goroutine (synchronous).

---

My intended purpose is to use this with https://github.com/WolvenSpirit/postgres-queue but it is just a few lines of code to provide a frame or example for any sort of node based separation of steps that together can define a task flow.

--- 

### Why define a generic task or a handler like this? 

Main reasons:

&#x2705; Handle errors consistently.

&#x2705; Handle program conditional logic flows consistently.

&#x2705; Have a recorded and traceable logic flow throughout task execution.

--- 

TODO:

&#x2705; More tests.

&#x2705; Now makes full use of Go Generics (go 1.18)

&#x10102; Provide clean-up (like `func onFinish() { clean(*n) }` ) hooks on each node. No matter if the node fails or succeeds, a custom clean-up function can run 
that might perform additional checks or store metrics regarding the outcome.

---

Examples for using this package can be found in the examples folder.
