package nodegraphflow

import "time"

// CircularNodePolicy enforces a shutdown procedure for the chain iteration
type CircularNodePolicy struct {
	Timeout            time.Duration
	RequiredForSuccess bool
	RestartOnError     bool
	IsCircularNode     bool
	StopChain          *chan int
}

// BuildChain will bind nodes into a chain which just means that the last node will have the first node as a sibling.
// Does not bind any siblings and no siblings should be set manually.
// The usecase for this would be running an refresh loop that has a lifetime which usually would equal that of the program.
func BuildChain(stopChain *chan int, nodes ...*Node) {
	for k := range nodes {
		if len(nodes)-1 != k {
			nodes[k].SubNodes = []*Node{nodes[k+1]}
		}
		if len(nodes)-1 == k {
			nodes[k].SubNodes = []*Node{nodes[0]}
		}
		nodes[k].CircularNodePolicy.IsCircularNode = true
		nodes[k].CircularNodePolicy.StopChain = stopChain
	}
}
