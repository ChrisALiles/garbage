package garbage

import (
	"fmt"
	"sync"
)

// Graph is the memory used by the mutator and managed by the
// collector.
// Node 0 is reserved for NIL.
// Node 1 is the head of the free list.
// Followed by the numRoots root nodes, then the rest.
var graph = make([]node, numNodes)

func initGraph() {
	// Populate memory linked to the root nodes.
	li := 2 + numRoots
	ri := numNodes - 1
	for i := 0; i < numRoots; i++ {
		prevli := 2 + i // start pointing at root
		for j := 0; j < 10; j++ {
			graph[prevli].leftEdge = li
			graph[li].rightEdge = ri
			prevli = li
			li += 1
			ri -= 1
		}
	}
	// Set up the free node list, initially containing all remaining
	// nodes linked by left edges.
	j := 1
	for i := numRoots*11 + 2; i < numNodes-numRoots*10; i++ {
		graph[j].leftEdge = i
		j = i
	}
}

// The mutexes prevent concurrent access to nodes.
// I use a number of them to restrict the range of mutual
// exclusion.
var mxs [numMxs]sync.Mutex

// Special mutex for the free list, for convenience.
var freeMx sync.Mutex

func lockNode(nodeNum int) {
	mxs[nodeNum%numMxs].Lock()
}
func unlockNode(nodeNum int) {
	mxs[nodeNum%numMxs].Unlock()
}
func shadeNode(nodeNum int) {
	if graph[nodeNum].colour == white {
		graph[nodeNum].colour = grey
	}
}

var freeChan = make(chan any)

// Unlink a node from the free list.
// If there are no free nodes, signal to the collector.
// If the collector can't help, return an error.
func getFreeNode() (int, error) {
	// First check if these is a free node - if not, signal.
	freeMx.Lock()
	if graph[1].leftEdge == 0 {
		freeMx.Unlock()
		freeChan <- true
		freeMx.Lock()
	}
	exFreeNode := graph[1].leftEdge
	if exFreeNode == 0 {
		freeMx.Unlock()
		return exFreeNode, fmt.Errorf("no free nodes")
	}
	graph[1].leftEdge = graph[exFreeNode].leftEdge
	freeMx.Unlock()
	return exFreeNode, nil
}
