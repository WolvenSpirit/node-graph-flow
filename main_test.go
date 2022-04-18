package nodegraphflow

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

type TestPayload struct{ State string }

func TestBindNodes(t *testing.T) {
	n1 := Node[TestPayload]{Name: "node1", Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("Node1")
		return TestPayload{}, nil
	}}
	n2 := Node[TestPayload]{Name: "node2", Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("Node2")
		return TestPayload{}, nil
	}}
	n3 := Node[TestPayload]{Name: "node3", Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("Node3")
		return TestPayload{}, nil
	}}
	type args struct {
		parent   *Node[TestPayload]
		siblings []*Node[TestPayload]
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "BindNodes test",
			args: args{parent: &n1, siblings: []*Node[TestPayload]{&n2, &n3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BindNodes(tt.args.parent, tt.args.siblings...)
		})
	}
}

func TestFlow(t *testing.T) {
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

	nA := Node[TestPayload]{Name: "nodeA", Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("NodeA")
		return TestPayload{}, nil
	}}
	nB := Node[TestPayload]{Name: "nodeB", ParentNode: &n1, Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("NodeB")
		return TestPayload{}, nil
	}}
	stopChain := make(chan int, 1)

	BindNodes(&n1, &n2)
	BindNodes(&n2, &n4, &n3)
	BindNodes(&n4, &n5)
	BindNodes(&n3, &n5)

	BuildChain(&stopChain, &nA, &nB)

	type args struct {
		n                *Node[TestPayload]
		i                TestPayload
		SubNodeIndex     int
		LateralNodeIndex int
		err              error
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Flow test",
			args: args{n: &n1, i: TestPayload{}, SubNodeIndex: 0, LateralNodeIndex: 0, err: nil}},
		{name: "Flow test",
			args: args{n: &n2, i: TestPayload{}, SubNodeIndex: 0, LateralNodeIndex: 0, err: errors.New("failed")}},
		{name: "Flow circular test",
			args: args{n: &nA, i: TestPayload{}, SubNodeIndex: 0, LateralNodeIndex: 0, err: nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				if tt.args.n.CircularNodePolicy.IsCircularNode {
					time.Sleep(time.Second * 1)
					(*tt.args.n.CircularNodePolicy.StopChain) <- 1
				}
			}()
			ctx := FlowContext{}
			ctx.Init()
			Flow(&ctx, tt.args.n, tt.args.i, tt.args.SubNodeIndex, tt.args.LateralNodeIndex, tt.args.err)

			if n5.Output.State != "Success" {
				t.Log("FlowTrail", n5.FlowTrail)
				t.Errorf("Failed to retrieve node output %+v", n5.Output)
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
	n1 := Node[TestPayload]{Name: "node1", Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("Node1")
		return TestPayload{}, nil
	}}
	n2 := Node[TestPayload]{Name: "node2", ParentNode: &n1, Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("Node2")
		return TestPayload{}, nil
	}}
	n3 := Node[TestPayload]{Name: "node3", ParentNode: &n2, Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("Node3")
		return TestPayload{}, nil
	}}
	n4 := Node[TestPayload]{Name: "node4", ParentNode: &n2, Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
		fmt.Println("Node4")
		return TestPayload{}, AbortError{}
	}}
	n5 := Node[TestPayload]{Name: "node5", ParentNode: &n3, Task: func(ctx *FlowContext, i TestPayload) (TestPayload, error) {
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
		n   *Node[TestPayload]
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Flow test",
			args: args{n: &n1, ctx: &FlowContext{Ctx: ctx, Cancel: cancel}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartFlow(tt.args.ctx, tt.args.n)
			if n4.FlowTrail == nil {
				t.Errorf("%+v", n4)
			}
			if n4.NodeTrail.NodeName == "" || n4.NodeTrail.NodeError == nil {
				t.Errorf("NodeTrail should be populated if error is returned")
			}
			if n2.NodeTrail.NodeName == "" || n2.NodeTrail.FinishedAt.Before(time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)) {
				t.Errorf("NodeTrail should be populated")
			}
			if n1.NodeTrail.NodeName == "" || n1.NodeTrail.FinishedAt.Before(time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)) {
				t.Errorf("NodeTrail should be populated")
			}
		})
	}
}
