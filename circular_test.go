package nodegraphflow

import (
	"fmt"
	"testing"
)

func TestBuildChain(t *testing.T) {
	n1 := Node{Name: "node1", Task: func(ctx *FlowContext, i Input) (Output, error) { fmt.Println("Node1"); return nil, nil }}
	n2 := Node{Name: "node2", ParentNode: &n1, Task: func(ctx *FlowContext, i Input) (Output, error) { fmt.Println("Node2"); return nil, nil }}
	n3 := Node{Name: "node3", Task: func(ctx *FlowContext, i Input) (Output, error) { fmt.Println("Node3"); return nil, nil }}
	n4 := Node{Name: "node4", Task: func(ctx *FlowContext, i Input) (Output, error) {
		fmt.Println("Node4")
		return nil, AbortError{}
	}}
	n5 := Node{Name: "node5", Task: func(ctx *FlowContext, i Input) (Output, error) {
		fmt.Println("Node5")
		return TestPayload{State: "Success"}, nil
	}}
	type args struct {
		nodes []*Node
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Test BuildChain", args: args{nodes: []*Node{&n1, &n2, &n3, &n4, &n5}}},
	}
	for _, tt := range tests {
		stopChain := make(chan int, 1)
		t.Run(tt.name, func(t *testing.T) {
			BuildChain(&stopChain, tt.args.nodes...)
		})
	}
}
