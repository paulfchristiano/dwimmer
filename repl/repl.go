package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/paulfchristiano/dwimmer"
)

var (
	cpuprofile = flag.String("cpu", "", "write cpu profile to file")
	memprofile = flag.String("mem", "", "write mem profile to file")
	state      = flag.String("state", "state", "load and save state to the DB entry with this key")
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(f)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	d := dwimmer.NewDwimmer(*state)
	defer d.Close()
	s := dwimmer.StartShell(d)
	if *memprofile != "" {
		runtime.GC()
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(f)
		}
		pprof.WriteHeapProfile(f)
		defer f.Close()
	}
	fmt.Println(s.Head())
}
