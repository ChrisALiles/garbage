package garbage

import (
	"fmt"
)

// The actual garbage collector.
// This was transcribed as closely as I could from the Dijkstra et al paper.
func collector(runTime int, colStats chan Stats,
	colTimeOut chan bool) {
	var (
		mi              int
		mk              int
		countPasses     int
		countFrees      int
		countFreeRelief int
	)
	for {
		// Shade the roots.
		for i := 0; i < numRoots; i++ {
			shadeNode(2 + i)
		}
		// Marking phase.
		mi = 2
		mk = numNodes
		// Marking cycle.
		for mk > 1 {
			left := 0
			right := 0
			lockNode(mi)
			if graph[mi].colour == grey {
				left = graph[mi].leftEdge
				right = graph[mi].rightEdge
				graph[mi].colour = black
				mk = numNodes
			} else {
				mk--
			}
			unlockNode(mi)
			// Shade the successors, if any.
			if left != 0 {
				lockNode(left)
				shadeNode(left)
				unlockNode(left)
			}
			if right != 0 {
				lockNode(right)
				shadeNode(right)
				unlockNode(right)
			}
			mi = (mi + 1) % numNodes
		}
		// Appending phase.
		mi = 2
		// Appending cycle.
		for mi < numNodes {
			lockNode(mi)
			if graph[mi].colour == white {
				// White mens garbge - append to the free list.
				graph[mi].rightEdge = 0
				freeMx.Lock()
				graph[mi].leftEdge = graph[1].leftEdge
				graph[1].leftEdge = mi
				freeMx.Unlock()
				unlockNode(mi)
				countFrees++
				// Check if the mutator is waiting for a free node.
				select {
				case <-freeChan:
					countFreeRelief++
				default:
				}
			} else {
				graph[mi].colour = white
				unlockNode(mi)
			}
			mi++
		}
		countPasses++
		// If there is a free node requesst here, then we are
		// out of free nodes altogether.
		// Let the mutator continue.
		select {
		case <-freeChan:
		default:
		}
		// Break and exit after notification from the mutator.
		select {
		case <-colTimeOut:
			fmt.Println("collector time out")
		default:
			continue
		}
		break
	}
	var st Stats
	st.NumFreeRelief = countFreeRelief
	st.NumFrees = countFrees
	st.NumPasses = countPasses
	colStats <- st
}
