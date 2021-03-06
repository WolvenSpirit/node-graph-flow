package nodegraphflow

import (
	"context"
	"os"
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

type Node[T interface{}] struct {
	Name               string                           // Name of the node
	ParentNode         *Node[T]                         // Parent node
	SubNodes           []*Node[T]                       // Children nodes
	Siblings           []*Node[T]                       // Lateral nodes
	Task               func(*FlowContext, T) (T, error) // Task that should be processed
	Input              T                                // Input payload, nil if starting node
	Output             T                                // Output payload
	FlowTrail          []string                         // The order in which nodes were executed
	NodeTrail          NodeTrail                        // Meta data populated after node processing finishes
	Context            *FlowContext                     // Pointer to the flow context
	CircularNodePolicy CircularNodePolicy               // Policy for circular nodes
	OnFinished         func(n T)                        // Called after task finishes execution
}

func (n *Node[Output]) SetOutput(o Output) {
	n.Output = o
}

func (n *Node[Input]) SetInput(i Input) {
	n.Input = i
}

// BindNodes links the parent to the sub nodes and each sub node laterally to each other.
func BindNodes[T interface{}](parent *Node[T], siblings ...*Node[T]) {
	parent.SubNodes = siblings
	for k := range siblings {
		siblings[k].Siblings = siblings
		siblings[k].ParentNode = parent
	}
}

// Flow initiates each node sequentially from the start node downstream to all sub nodes.
// If a parent has more than one sub node, the higher index nodes are fallback nodes.
// Should the first of the siblings fail, the next lateral node will execute from the siblings slice.
// If all nodes from a level error out then the context of the flow will be canceled.
func Flow[T interface{}](ctx *FlowContext, n *Node[T], i T, SubNodeIndex int, LateralNodeIndex int, err error) {
	if n.CircularNodePolicy.StopChain != nil {
		select {
		case <-*n.CircularNodePolicy.StopChain:
			return
		default:
			// Just run
		}
	}
	if t, errctx := ctx.IsCanceled(); t {
		if errctx != nil {
			os.Stderr.Write([]byte(errctx.Error()))
		}
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
	if n.OnFinished != nil {
		n.OnFinished(o)
	}
	n.SetOutput(o)
	nt.FinishedAt = time.Now()
	nt.NodeError = err
	n.FlowTrail = []string{n.Name}
	n.NodeTrail = nt
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
func StartFlow[T interface{}](ctx *FlowContext, n *Node[T]) {
	var e T
	Flow(ctx, n, e, 0, 0, nil)
}
