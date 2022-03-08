package nodegraphflow

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

type TestPayload struct{ State string }

func TestBindNodes(t *testing.T) {
	n1 := Node{Name: "node1", Task: func(ctx *FlowContext, i Input) (Output, error) { fmt.Println("Node1"); return nil, nil }}
	n2 := Node{Name: "node2", Task: func(ctx *FlowContext, i Input) (Output, error) { fmt.Println("Node2"); return nil, nil }}
	n3 := Node{Name: "node3", Task: func(ctx *FlowContext, i Input) (Output, error) { fmt.Println("Node3"); return nil, nil }}
	type args struct {
		parent   *Node
		siblings []*Node
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "BindNodes test",
			args: args{parent: &n1, siblings: []*Node{&n2, &n3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BindNodes(tt.args.parent, tt.args.siblings...)
		})
	}
}

func TestFlow(t *testing.T) {
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

	BindNodes(&n1, &n2)
	BindNodes(&n2, &n4, &n3)
	BindNodes(&n4, &n5)
	BindNodes(&n3, &n5)

	type args struct {
		n                *Node
		i                Input
		SubNodeIndex     int
		LateralNodeIndex int
		err              error
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Flow test",
			args: args{n: &n1, i: nil, SubNodeIndex: 0, LateralNodeIndex: 0, err: nil}},
		{name: "Flow test",
			args: args{n: &n2, i: nil, SubNodeIndex: 0, LateralNodeIndex: 0, err: errors.New("failed")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := FlowContext{}
			ctx.Init()
			Flow(&ctx, tt.args.n, tt.args.i, tt.args.SubNodeIndex, tt.args.LateralNodeIndex, tt.args.err)
			if o, ok := n5.Output.(TestPayload); ok {
				if o.State != "Success" {
					t.Log("FlowTrail", n5.FlowTrail)
					t.Errorf("Failed to retrieve node output %+v", o)
				}
			}
		})
	}
}

func TestAbortError_Error(t *testing.T) {
	type fields struct {
		Message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "AbortError test",
			fields: fields{Message: "error"},
			want:   "Flow through nodes has been aborted",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AbortError{
				Message: tt.fields.Message,
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("AbortError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartFlow(t *testing.T) {
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

	BindNodes(&n1, &n2)
	BindNodes(&n2, &n4, &n3)
	BindNodes(&n4, &n5)
	BindNodes(&n3, &n5)
	ctx, cancel := context.WithCancel(context.Background())
	type args struct {
		ctx *FlowContext
		n   *Node
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Flow test",
			args: args{n: &n1, ctx: &FlowContext{ctx: ctx, cancel: cancel}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartFlow(tt.args.ctx, tt.args.n)
		})
	}
}
