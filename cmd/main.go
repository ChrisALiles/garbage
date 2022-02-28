package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/ChrisALiles/garbage"
)

var (
	runTime  int
	wantProf bool
)

func main() {
	getflags()

	if wantProf {
		f, err := os.Create("./gprof")
		if err != nil {
			panic("cannot create profile file")
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	// Turn on interrupt signal handler for debugging.
	go garbage.WaitForSignal()

	// Start the garbage collection simulation and print the stats.
	statChan := make(chan garbage.Stats)
	go garbage.Run(runTime, statChan)
	fmt.Println(<-statChan)
}

func getflags() {
	// Parse command line flags.
	rt := flag.Int("t", 1, "Run time in seconds")
	pr := flag.Bool("p", false, "Turn on CPU profiling")

	flag.Parse()
	runTime = *rt
	wantProf = *pr
}
