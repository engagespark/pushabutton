package pushabutton

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const (
	assetsDir             = "assets"
	buttonsDir            = "buttons"
	logsDir               = "logs"
	exampleScriptFilename = "what_is_the_current_date.sh"
)

var logfilePath = path.Join(logsDir, "journal.log")

func Setup() {
	vanilla := !FileExists(buttonsDir)
	if vanilla {
		fmt.Println("Buttons directory missing. Creating it for you and filling it with examples.")
		createButtonLibDir()
		createExampleScript()
	} else {
		fmt.Println("Buttons dir exists, skipping.")
	}
	if !FileExists(logsDir) {
		fmt.Println("Logs directory missing. Creating it for you.")
		createLogsDir()
	} else {
		fmt.Println("Logs dir exists, skipping.")
	}
}

func createButtonLibDir() {
	if err := os.Mkdir(buttonsDir, 0700); err != nil {
		fmt.Errorf("Could not create directory for buttons: ", err)
	}
	fmt.Printf("Created buttons directory: ./%v\n", buttonsDir)

}

func createLogsDir() {
	if err := os.Mkdir(logsDir, 0700); err != nil {
		fmt.Errorf("Could not create directory for logs: ", err)
	}
	fmt.Printf("Created logs directory: ./%v\n", logsDir)

}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	fmt.Println(err)
	return true
}

func createExampleScript() {
	fmt.Println("Checking example script.")
	targetPath := path.Join(buttonsDir, exampleScriptFilename)
	sourcePath := path.Join(assetsDir, exampleScriptFilename)
	if FileExists(targetPath) {
		fmt.Println("Example script exists, not touching it: ", targetPath)
		return
	}
	data, err := Asset(sourcePath)
	if err != nil {
		fmt.Printf("Failed generating the example script, sorry: %v", err)
		return
	}
	fmt.Println("Writing example script: ", targetPath)
	ioutil.WriteFile(targetPath, data, 0700)
}
