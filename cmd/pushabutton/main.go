package main

import (
	"fmt"
	"os"

	"github.com/mknecht/pushabutton"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app   = kingpin.New("pushabutton", "A web application executing your scripts.")
	debug = app.Flag("debug", "enable debug mode").Default("false").Bool()
	setup = app.Command("setup", "Setup vanilla config or repair missing essentials. Operates on working directory.")
	serve = app.Command("serve", "Run a webserver, waiting to run scripts.").Default()
	addr  = serve.Arg("addr", "Where the webserver should listen.").Default(":8080").String()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.CommandLine.HelpFlag.Short('h')
	app.UsageTemplate(kingpin.DefaultUsageTemplate)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case setup.FullCommand():
		fmt.Println("Setting up vanilla config, or repairing it …")
		pushabutton.Setup()
		break
	case serve.FullCommand():
		fmt.Println("Running server")
		pushabutton.StartServerOrCrash(*addr)
		break
	}
}
