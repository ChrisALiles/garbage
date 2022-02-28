package garbage

import (
	"fmt"
	"math/rand"
	"time"
)

// The mutator does some memory manipulation -
// 1) Remove a node from the free list and attach it to active memory
// 2) Produce grabage by detaching a mode from active memory
// and also tries to mimic actual work by repeated passes through
// ctive memory.
func mutator(runTime int, colTimeOut chan bool) {
	timeout := time.After(time.Duration(runTime * int(time.Second)))
	for {
		// Find a node to attach the free one to.
		root := rand.Intn(numRoots)
		nodeNum := graph[root+2].leftEdge
		for graph[nodeNum].leftEdge != 0 {
			nodeNum = graph[nodeNum].leftEdge
		}
		// Link the free node and remove it from the free list.
		exFreeNode, err := getFreeNode()
		// Error returned if no free nodes available.
		// we just skip it for this pass.
		if err == nil {
			lockNode(nodeNum)
			graph[nodeNum].leftEdge = exFreeNode
			unlockNode(nodeNum)
			lockNode(exFreeNode)
			shadeNode(exFreeNode)
			graph[exFreeNode].leftEdge = 0
			graph[exFreeNode].rightEdge = 0
			unlockNode(exFreeNode)
		}
		// Now do some "work"
		var dummySum int
		for i := 0; i < 1000000; i++ {
			for j := 2; j < numRoots+2; j++ {
				k := j
				for graph[k].leftEdge != 0 {
					dummySum += graph[k].leftEdge
					k = graph[k].leftEdge
				}
			}
		}
		// Create garbage by removing the edges from a node.
		root = rand.Intn(numRoots)
		nodeNum = graph[root+2].leftEdge
		prevNodeNum := 0
		for graph[nodeNum].leftEdge != 0 {
			prevNodeNum = nodeNum
			nodeNum = graph[nodeNum].leftEdge
		}
		if prevNodeNum != 0 {
			lockNode(prevNodeNum)
			graph[prevNodeNum].leftEdge = 0
			graph[prevNodeNum].rightEdge = 0
			unlockNode(prevNodeNum)
		}
		// Now do some more "work"
		for i := 0; i < 1000000; i++ {
			for j := 2; j < numRoots+2; j++ {
				k := j
				for graph[k].rightEdge != 0 {
					dummySum += graph[k].rightEdge
					k = graph[k].rightEdge
				}
			}
		}
		// Break and exit after the timeout.
		select {
		case <-timeout:
			fmt.Println("mutator time out")
			// Tell the collector to stop.
			colTimeOut <- true
		default:
			continue
		}
		break
	}

}
