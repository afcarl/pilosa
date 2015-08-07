package main

import (
	"flag"
	"os"
	"runtime/pprof"

	log "github.com/cihub/seelog"
	"github.com/mitchellh/panicwrap"
	"github.com/umbel/pilosa/config"
	"github.com/umbel/pilosa/core"
	"github.com/umbel/pilosa/cruncher"
	"github.com/umbel/pilosa/index"
	"github.com/umbel/pilosa/util"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	Build      string
)

func main() {
	defer log.Flush()
	exitStatus, err := panicwrap.BasicWrap(panicHandler)
	if err != nil {
		// Something went wrong setting up the panic wrapper. Unlikely,
		// but possible.
		panic(err)
	}

	if exitStatus >= 0 {
		os.Exit(exitStatus)
	}
	core.Build = Build

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Warn(err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	config.SetupConfig()
	util.SetupUtil()
	index.SetupCassandra()

	cruncher := cruncher.NewCruncher()
	cruncher.Run()
	log.Warn("STOP")
}

func panicHandler(output string) {
	// output contains the full output (including stack traces) of the
	// panic. Put it in a file or something.
	log.Critical("The child panicked:\n\n", output)
	os.Exit(1)
}