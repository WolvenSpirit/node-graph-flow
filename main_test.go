package nodegraphflow

import (
	"errors"
	"fmt"
	"testing"
)

func TestBindNodes(t *testing.T) {
	n1 := Node{Name: "node1", Task: func(i Input) (Output, error) { fmt.Println("Node1"); return nil, nil }}
	n2 := Node{Name: "node2", Task: func(i Input) (Output, error) { fmt.Println("Node2"); return nil, nil }}
	n3 := Node{Name: "node3", Task: func(i Input) (Output, error) { fmt.Println("Node3"); return nil, nil }}
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
	n1 := Node{Name: "node1", Task: func(i Input) (Output, error) { fmt.Println("Node1"); return nil, nil }}
	n2 := Node{Name: "node2", Task: func(i Input) (Output, error) { fmt.Println("Node2"); return nil, nil }}
	n3 := Node{Name: "node3", Task: func(i Input) (Output, error) { fmt.Println("Node3"); return nil, nil }}
	n4 := Node{Name: "node4", Task: func(i Input) (Output, error) { fmt.Println("Node4"); return nil, errors.New("failed") }}
	n5 := Node{Name: "node5", Task: func(i Input) (Output, error) { fmt.Println("Node5"); return nil, nil }}

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
			Flow(tt.args.n, tt.args.i, tt.args.SubNodeIndex, tt.args.LateralNodeIndex, tt.args.err)
		})
	}
}
