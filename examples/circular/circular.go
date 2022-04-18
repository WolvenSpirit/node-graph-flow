package main

import (
	"context"
	"time"

	flow "github.com/WolvenSpirit/node-graph-flow"
)

type TestPayload struct{ State string }

func main() {
	stopChain := make(chan int, 1)

	node1 := flow.Node[TestPayload]{Task: func(fc *flow.FlowContext, i TestPayload) (TestPayload, error) {
		println("node1 executing")
		return TestPayload{}, nil
	}, Name: "node1"}
	node2 := flow.Node[TestPayload]{Task: func(fc *flow.FlowContext, i TestPayload) (TestPayload, error) {
		println("node2 executing")
		return TestPayload{}, nil
	}, Name: "node2"}
	node3 := flow.Node[TestPayload]{Task: func(fc *flow.FlowContext, i TestPayload) (TestPayload, error) {
		println("node3 executing")
		return TestPayload{}, nil
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
