package nodegraphflow

type Input interface{}
type Output interface{}

/*
type FlowContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *FlowContext) Init() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
}

func (c *FlowContext) IsCanceled() (bool, error) {
	if err := c.ctx.Err(); err != nil {
		return true, err
	}
	return false, nil
}
*/

type Tracker struct {
	CurrentStreamNodeDepth int
	CurrentLateralNode     int
}

type Node struct {
	Name       string                      // Name of the node
	ParentNode *Node                       // Parent node
	SubNodes   []*Node                     // Children nodes
	Siblings   []*Node                     // Lateral nodes
	Task       func(Input) (Output, error) // Task that should be processed
	Input      Input                       // Input payload, nil if starting node
	Output     Output                      // Output payload
}

// BindNodes links the parent to the sub nodes and each sub node laterally to each other.
func BindNodes(parent *Node, siblings ...*Node) {
	parent.SubNodes = siblings
	for k := range siblings {
		siblings[k].Siblings = siblings
	}
}

// Flow initiates each node sequentially from the start node downstream to all sub nodes.
// If a parent has more than one sub node, the higher index nodes are fallback nodes.
// Should the first of the siblings fail, the next lateral node will execute from the siblings slice.
// If all nodes from a level error out then the context of the flow will be canceled (TODO).
func Flow(n *Node, i Input, SubNodeIndex int, LateralNodeIndex int, err error) {
	if err != nil {
		Flow(n.Siblings[LateralNodeIndex], i, SubNodeIndex, LateralNodeIndex, nil)
	}
	o, err := n.Task(i)
	if len(n.SubNodes) != 0 && err == nil {
		Flow(n.SubNodes[SubNodeIndex], o, SubNodeIndex, LateralNodeIndex, err)
	}
	if err != nil {
		LateralNodeIndex++
		Flow(n.Siblings[LateralNodeIndex], i, SubNodeIndex, LateralNodeIndex, nil)
	}
}

/*
func main() {

	n1 := Node{Name: "node1", Task: func(i Input) (Output, error) { fmt.Println("Node1"); return nil, nil }}
	n2 := Node{Name: "node2", Task: func(i Input) (Output, error) { fmt.Println("Node2"); return nil, nil }}
	n3 := Node{Name: "node3", Task: func(i Input) (Output, error) { fmt.Println("Node3"); return nil, nil }}
	n4 := Node{Name: "node4", Task: func(i Input) (Output, error) { fmt.Println("Node4"); return nil, errors.New("failed") }}
	n5 := Node{Name: "node5", Task: func(i Input) (Output, error) { fmt.Println("Node5"); return nil, nil }}

	BindNodes(&n1, &n2)
	BindNodes(&n2, &n4, &n3)
	BindNodes(&n4, &n5)
	BindNodes(&n3, &n5)
	Flow(&n1, nil, 0, 0, nil)
}
*/
