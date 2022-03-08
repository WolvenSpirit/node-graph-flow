package main

import (
	"context"
	"time"

	flow "github.com/WolvenSpirit/node-graph-flow"
)

func main() {
	stopChain := make(chan int, 1)

	node1 := flow.Node{Task: func(fc *flow.FlowContext, i flow.Input) (flow.Output, error) {
		println("node1 executing")
		return nil, nil
	}, Name: "node1"}
	node2 := flow.Node{Task: func(fc *flow.FlowContext, i flow.Input) (flow.Output, error) {
		println("node2 executing")
		return nil, nil
	}, Name: "node2"}
	node3 := flow.Node{Task: func(fc *flow.FlowContext, i flow.Input) (flow.Output, error) {
		println("node3 executing")
		return nil, nil
	}, Name: "node3"}

	flow.BuildChain(&stopChain, &node1, &node2, &node3)

	ctx, cancel := context.WithCancel(context.Background())
	fctx := flow.FlowContext{Ctx: ctx, Cancel: cancel}

	go func() {
		println("Node chain loop start")
		flow.StartFlow(&fctx, &node1)
		println("Chain loop stop")
	}()

	time.Sleep(time.Second * 3)
	stopChain <- 1
}
