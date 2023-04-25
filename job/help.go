package job

import (
	"fmt"
	"os"
	"strings"
)

type helpJob struct{}

func (j *helpJob) Run() error {
	help := `
Usage: vault2env [--config=<config_file>] [--from=<uri>] [--to=<uri>] [preset]

  A <uri> is either a file://path or vault://path
  A [preset] in your config file contains pre-defined *from* and *to* values
  --from and --to always take priority over the preset
  Piping always takes priority over any other provided options

Examples:

  vault2env --from=file://.env --to=vault://secret/your/secret/path
  vault2env pull # Executes a preset
  vault2env --to=file://.other.env pull # override *to* from preset
  cat .test.env | vault2env push
  vault2env pull > .some.env
`

	fmt.Fprintln(os.Stderr, strings.TrimSpace(help))

	return nil
}
