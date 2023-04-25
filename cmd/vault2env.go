package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rmarganti/vault2env/config"
	"github.com/rmarganti/vault2env/job"
)

func main() {
	args := parseArgs()

	cfg, err := config.Load(args.configPath)
	checkErr(err)

	job, err := job.New(cfg, args.from, args.to, args.preset)
	checkErr(err)

	err = job.Run()
	checkErr(err)
}

type Args struct {
	configPath string
	from       string
	to         string
	preset     string
}

func parseArgs() Args {
	configPath := flag.String("config", config.DefaultFilePath, "Config file path")
	from := flag.String("from", "", "Source to pull from")
	to := flag.String("to", "", "Source to write to")
	flag.Parse()

	preset := ""
	if len(flag.Args()) > 0 {
		preset = flag.Args()[0]
	}

	return Args{
		configPath: *configPath,
		from:       *from,
		to:         *to,
		preset:     preset,
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
