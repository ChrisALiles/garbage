package garbage

import (
	"fmt"
)

// Unit of memory.
type node struct {
	leftEdge  int
	rightEdge int
	colour    uint8
}

// Stats gathered by the collector.
type Stats struct {
	NumPasses     int
	NumFrees      int
	NumFreeRelief int
}

func (s Stats) String() string {
	return fmt.Sprintln("Collector passes:", s.NumPasses,
		"\nCollector nodes freed:", s.NumFrees,
		"\nCollector free reliefs:", s.NumFreeRelief)
}
