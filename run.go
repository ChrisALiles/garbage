// Garbage is an exploration of the garbage collection technique
// described in http://dl.acm.org/citation.cfm?id=359655.
// Dijkstra's method is used by Go's garbage collector.
package garbage

import (
	"sync"
)

// Start the collector and mutator processes.
// the mutator is terminated after a time after timeoout,
// at which point it notifies the collector it's time to
// finish, via the colTimeOut channel.
// The collector can't be allowed to finish first becuase
// the mutator might then get hung on a free node request.
func Run(runTime int, statChan chan Stats) {
	var wg sync.WaitGroup
	var colStats = make(chan Stats)
	var colTimeOut = make(chan bool, 2)
	wg.Add(1)
	go func() {
		defer wg.Done()
		collector(runTime, colStats, colTimeOut)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		mutator(runTime, colTimeOut)
	}()
	st := <-colStats
	wg.Wait()
	statChan <- st
}
