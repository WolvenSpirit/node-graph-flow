package nodegraphflow

import (
	"fmt"
	"testing"
)

func TestBuildChain(t *testing.T) {
	n1 := Node[TestPayload]{Name: "node1", Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("Node1")
		return TestPayload{}, nil
	}}
	n2 := Node[TestPayload]{Name: "node2", ParentNode: &n1, Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("Node2")
		return TestPayload{}, nil
	}}
	n3 := Node[TestPayload]{Name: "node3", Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("Node3")
		return TestPayload{}, nil
	}}
	n4 := Node[TestPayload]{Name: "node4", Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("Node4")
		return TestPayload{}, AbortError{}
	}}
	n5 := Node[TestPayload]{Name: "node5", Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("Node5")
		return TestPayload{State: "Success"}, nil
	}}
	type args struct {
		nodes []*Node[TestPayload]
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Test BuildChain", args: args{nodes: []*Node[TestPayload]{&n1, &n2, &n3, &n4, &n5}}},
	}
	for _, tt := range tests {
		stopChain := make(chan int, 1)
		t.Run(tt.name, func(t *testing.T) {
			BuildChain(&stopChain, tt.args.nodes...)
		})
	}
}
