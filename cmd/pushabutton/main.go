package main

import (
	"fmt"
	"os"

	"github.com/mknecht/pushabutton"
	"gopkg.in/alecthomas/kingpin.v2"
)

// compile passing -ldflags "-X main.Build <build sha1>"
var Build string

var (
	app     = kingpin.New("pushabutton", "A web application executing your scripts.")
	setup   = app.Command("setup", "Setup vanilla config or repair missing essentials. Operates on working directory.")
	serve   = app.Command("serve", "Run a webserver, waiting to run scripts.").Default()
	addr    = serve.Arg("addr", "Where the webserver should listen.").Default(":8080").String()
	baseUrl = serve.Flag("base-url", "The base URL for the webapp.").Default("/").String()
)

func main() {
	kingpin.CommandLine.HelpFlag.Short('h')
	app.Version(Build)
	app.UsageTemplate(kingpin.DefaultUsageTemplate)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case setup.FullCommand():
		fmt.Println("Setting up vanilla config, or repairing it …")
		pushabutton.Setup()
		break
	case serve.FullCommand():
		fmt.Printf("Running server on %v\n", *addr)
		pushabutton.StartServerOrCrash(*addr, *baseUrl)
		break
	}
}
