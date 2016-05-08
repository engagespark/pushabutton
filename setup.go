package pushabutton

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const (
	assetsDir      = "assets"
	buttonsDir     = "buttons"
	scriptFilename = "what_is_the_current_date.sh"
)

func Setup() {
	vanilla := !exists(buttonsDir)
	if vanilla {
		fmt.Println("Buttons directory missing. Creating it for you and filling it with examples.")
		createButtonLibDir()
		createExampleScript()
	} else {
		fmt.Println("Buttons dir exists, skipping.")
	}
}

func createButtonLibDir() {
	if err := os.Mkdir(buttonsDir, 0700); err != nil {
		fmt.Errorf("Could not create directory for buttons: ", err)
	}
	fmt.Printf("Created buttons directory: ./%v\n", buttonsDir)

}

func exists(path string) bool {
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
	targetPath := path.Join(buttonsDir, scriptFilename)
	sourcePath := path.Join(assetsDir, scriptFilename)
	if exists(targetPath) {
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
