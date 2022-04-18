package main

import (
	"context"

	flow "github.com/WolvenSpirit/node-graph-flow"
)

type myPayload struct {
	Data string
}

func main() {
	node1 := flow.Node[myPayload]{Task: func(fc *flow.FlowContext, i myPayload) (myPayload, error) {
		return myPayload{Data: "node1 Output\n"}, nil
	}, Name: "node1"}
	node2 := flow.Node[myPayload]{Task: func(fc *flow.FlowContext, i myPayload) (myPayload, error) {
		i.Data += "node2 Output\n"
		return i, nil
	}, Name: "node2"}
	node3 := flow.Node[myPayload]{Task: func(fc *flow.FlowContext, i myPayload) (myPayload, error) {
		i.Data += "node3 Output\n"
		return i, nil
	}, Name: "node3"}
	flow.BindNodes(&node1, &node2)
	flow.BindNodes(&node2, &node3)
	ctx, cancel := context.WithCancel(context.Background())
	fctx := flow.FlowContext{Ctx: ctx, Cancel: cancel}
	flow.StartFlow(&fctx, &node1)
	println(node3.Output.Data)
}
