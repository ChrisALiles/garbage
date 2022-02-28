package garbage

import (
	"fmt"
	"os"
	"os/signal"
)

func WaitForSignal() {
	var sig os.Signal

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	sig = <-c
	fmt.Println("signal", sig)
	panic("^C received - where were we?")
}
