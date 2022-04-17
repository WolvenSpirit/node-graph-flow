[![node-graph-flow](https://github.com/WolvenSpirit/node-graph-flow/actions/workflows/go.yml/badge.svg)](https://github.com/WolvenSpirit/node-graph-flow/actions/workflows/go.yml) [![codecov](https://codecov.io/gh/WolvenSpirit/node-graph-flow/branch/main/graph/badge.svg?token=hDXMUdD4L1)](https://codecov.io/gh/WolvenSpirit/node-graph-flow)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/WolvenSpirit/node-graph-flow)

## node graph flow

Minimal frame for node graph based task processing.

### Disclaimer: 

- Each individual flow can be executed concurrently on separate goroutines.

- At node level, all nodes of an individual graph should execute on the same goroutine (synchronous).

---

My intended purpose is to use this with https://github.com/WolvenSpirit/postgres-queue but it is just a few lines of code to provide a frame or example for any sort of node based separation of steps that together can define a task flow.

--- 

### Why define a generic task or a handler like this? 

In Go, handling errors always creates branches within a program (like in most languages), this might make things hard to debug later on.

In contrast, what others consider good code, slim functions with a single purpose etc. does tend to just wrap to many things and makes 
another mess because things now at top level are just some custom named wrappers that might not even be in the same style within a program. 

This brings more confusion.

This package arguably is one way to solve the problems and provides error handling with errors that are recorded throughout the task execution.

--- 

TODO:

- More tests.

- Provide clean-up (like `func onFinish() { clean(*n) }` ) hooks on each node. No matter if the node fails or succeeds, a custom clean-up function can run 
that might perform additional checks or store metrics regarding the outcome.

---

Examples for using this package can be found in the examples folder.
