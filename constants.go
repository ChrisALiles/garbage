package garbage

const (
	white = iota // marking colours.
	grey
	black

	numNodes = 1000 // size of memory.
	numRoots = 10   // number of root nodes.
	numMxs   = 10   // number of mutexes to control memory access.
)
