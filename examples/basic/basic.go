package main

import (
	"context"

	flow "github.com/WolvenSpirit/node-graph-flow"
)

func main() {
	node1 := flow.Node{Task: func(fc *flow.FlowContext, i flow.Input) (flow.Output, error) {
		return "node1 Output\n", nil
	}, Name: "node1"}
	node2 := flow.Node{Task: func(fc *flow.FlowContext, i flow.Input) (flow.Output, error) {
		return i.(string) + "node2 Output\n", nil
	}, Name: "node2"}
	node3 := flow.Node{Task: func(fc *flow.FlowContext, i flow.Input) (flow.Output, error) {
		return i.(string) + "node3 Output\n", nil
	}, Name: "node3"}
	flow.BindNodes(&node1, &node2)
	flow.BindNodes(&node2, &node3)
	ctx, cancel := context.WithCancel(context.Background())
	fctx := flow.FlowContext{Ctx: ctx, Cancel: cancel}
	flow.StartFlow(&fctx, &node1)
	println(node3.Output.(string))
}
