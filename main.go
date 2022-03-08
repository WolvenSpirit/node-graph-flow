package nodegraphflow

import (
	"context"
	"time"
)

type Input interface{}
type Output interface{}

type AbortError struct {
	Message string
}

func (err AbortError) Error() string {
	return "Flow through nodes has been aborted"
}

// Initialized context and cancel func need to be populated
type FlowContext struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

func (c *FlowContext) Init() {
	c.Ctx, c.Cancel = context.WithCancel(context.Background())
}

func (c *FlowContext) IsCanceled() (bool, error) {
	if err := c.Ctx.Err(); err != nil {
		return true, err
	}
	return false, nil
}

type NodeTrail struct {
	NodeName   string
	StartedAt  time.Time
	FinishedAt time.Time
	NodeError  error
}

type Node struct {
	Name               string                                    // Name of the node
	ParentNode         *Node                                     // Parent node
	SubNodes           []*Node                                   // Children nodes
	Siblings           []*Node                                   // Lateral nodes
	Task               func(*FlowContext, Input) (Output, error) // Task that should be processed
	Input              Input                                     // Input payload, nil if starting node
	Output             Output                                    // Output payload
	FlowTrail          []string                                  // The order in which nodes were executed
	NodeTrail          NodeTrail                                 // Meta data populated after node processing finishes
	Context            *FlowContext                              // Pointer to the flow context
	CircularNodePolicy CircularNodePolicy                        // Policy for circular nodes
}

func (n *Node) SetOutput(o Output) {
	n.Output = o
}

func (n *Node) SetInput(i Input) {
	n.Input = i
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
func Flow(ctx *FlowContext, n *Node, i Input, SubNodeIndex int, LateralNodeIndex int, err error) {
	if n.CircularNodePolicy.StopChain != nil {
		select {
		case <-*n.CircularNodePolicy.StopChain:
			return
		default:
			// Just run
		}
	}
	// TODO properly check if it's canceled
	if t, _ := ctx.IsCanceled(); t {
		return
	}
	if err != nil {
		Flow(ctx, n.Siblings[LateralNodeIndex], i, SubNodeIndex, LateralNodeIndex, nil)
		return
	}
	nt := NodeTrail{NodeName: n.Name, StartedAt: time.Now()}
	n.SetInput(i)
	n.Context = ctx
	o, err := n.Task(ctx, i)
	n.SetOutput(o)
	nt.FinishedAt = time.Now()
	nt.NodeError = err
	n.FlowTrail = []string{n.Name}
	if n.ParentNode != nil {
		n.FlowTrail = append(n.ParentNode.FlowTrail, n.FlowTrail...)
	}
	if _, ok := err.(AbortError); ok {
		ctx.Cancel() // Unrecoverable error or inability to further process downstream with currently available inputs or no lateral nodes left
	}
	if len(n.SubNodes) != 0 && err == nil {
		Flow(ctx, n.SubNodes[SubNodeIndex], o, SubNodeIndex, LateralNodeIndex, err)
	}
	if err != nil {
		LateralNodeIndex++
		Flow(ctx, n.Siblings[LateralNodeIndex], i, SubNodeIndex, LateralNodeIndex, nil)
	}
}

// StartFlow is an alias for calling Flow with the arguments as Flow(ctx, n, nil, 0, 0, nil)
func StartFlow(ctx *FlowContext, n *Node) {
	Flow(ctx, n, nil, 0, 0, nil)
}
