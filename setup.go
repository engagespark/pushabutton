package pushabutton

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const (
	assetsDir  = "assets"
	buttonsDir = "buttons"
	logsDir    = "logs"
)

var exampleFiles = []string{
	"what_is_the_current_date.sh",
	"write-message-to-logged-in-user.sh",
	"write-message-to-logged-in-user.sh.parameters",
	"write-message-to-logged-in-user.sh.parameters.user-tty.choices.sh",
}

var logfilePath = path.Join(logsDir, "journal.log")

func Setup() {
	vanilla := !FileExists(buttonsDir)
	if vanilla {
		fmt.Println("Buttons directory missing. Creating it for you and filling it with examples.")
		createButtonLibDir()
		createExampleScripts()
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

func createExampleScripts() {
	for _, filename := range exampleFiles {
		fmt.Printf("Checking example file: %v\n", filename)
		targetPath := path.Join(buttonsDir, filename)
		sourcePath := path.Join(assetsDir, filename)
		if FileExists(targetPath) {
			fmt.Printf("Example script exists, not touching it: %v\n", targetPath)
			return
		}
		data, err := Asset(sourcePath)
		if err != nil {
			fmt.Printf("Failed generating the example script, sorry: %v\n", err)
			return
		}
		fmt.Printf("Writing example script: %v\n", targetPath)
		ioutil.WriteFile(targetPath, data, 0700)
	}
}
