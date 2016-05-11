package pushabutton

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/nu7hatch/gouuid"
)

const (
	parametersSuffix    = ".parameters"
	choicesSuffix       = ".choices.sh"
	parameterTypeString = "string"
	parameterTypeChoice = "choice"
)

var knownParameterTypes = []string{parameterTypeChoice, parameterTypeString}

func StartServerOrCrash(addr string) {
	http.Handle("/static/", http.StripPrefix("/static/", ServeAsset{}))
	http.Handle("/push/", http.StripPrefix("/push/", PostPush{}))
	http.Handle("/buttons", GetButtons{})
	http.Handle("/", ServeIndex{})
	log.Fatal(http.ListenAndServe(addr, nil))
}

type ServeAsset struct{}

func (handler ServeAsset) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.URL.Path)
	data, err := Asset(path.Join(assetsDir, filename))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "File not found.")
		fmt.Printf("File %v not found: %v\n", filename, err)
		return
	}

	w.Header().Set("Content-Type", "text/css")
	fmt.Fprint(w, string(data))

}

type ServeIndex struct{}

func (handler ServeIndex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := Asset(path.Join(assetsDir, "index.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Template not found")
		fmt.Printf("Template for page not found: %v\n", err)
		return
	}
	indexHtml := string(data)
	tmpl, err := template.New("index.html").Parse(indexHtml)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not parse template.")
		fmt.Printf("Could not parse template for page: %v\n", err)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not execute template.")
		fmt.Printf("Could not execute template for page: %v\n", err)
	}

	err = tmpl.Execute(w, map[string]string{})

}

type PostPush struct{}

func (handler PostPush) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buttonId := r.URL.Path

	pushId, err := pushButton(buttonId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not push button.")
		fmt.Printf("Could not push button: %v\n", err)
	}

	responseData := map[string]string{"buttonId": buttonId, "pushId": pushId}
	payload, err := json.Marshal(responseData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not encode payload as JSON.")
		fmt.Printf("Could not encode payload as JSON: %v\n", err)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(payload)
	fmt.Sprintf("Pressed %v: %v\n", buttonId, pushId)
}

type ParameterDef struct {
	Name        string
	Type        string
	Description string
	Details     map[string]interface{}
}

type Button struct {
	Id         string
	Title      string
	Parameters []ParameterDef
}

type GetButtons struct{}

func (handler GetButtons) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := json.Marshal(AvailableButtons())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not encode payload as JSON.")
		fmt.Printf("Could not encode payload as JSON: %v\n", err)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(payload)
}

func AvailableButtons() []Button {
	var buttons []Button
	filepath.Walk(buttonsDir, func(candidate string, info os.FileInfo, err error) error {
		if candidate == buttonsDir {
			return nil
		}
		if info.Mode()&os.ModePerm&0100 == 0 {
			fmt.Printf("skipping %v\n", candidate)
			return nil
		}

		filename := path.Base(candidate)

		if strings.TrimSpace(filename) == "" {
			fmt.Printf("skipping suspicous whitespace file %q\n", candidate)
			return nil
		}

		if len(buttons) > 0 && strings.HasPrefix(filename, buttons[len(buttons)-1].Id) {
			fmt.Printf("skipping parameter file %v of %v\n", candidate, buttons[len(buttons)-1].Id)
			return nil
		}

		buttonId := filename
		title := generateTitle(filename)
		parameters, err := loadParameters(filename)
		if err != nil {
			fmt.Printf("Error loading button from %v: %v\n", filename, err)
			return nil
		}

		buttons = append(buttons, Button{
			Id:         html.EscapeString(buttonId),
			Title:      html.EscapeString(title),
			Parameters: parameters,
		})
		return nil
	})
	return buttons
}

func generateTitle(filename string) string {
	questionWords := []string{"how", "what", "who", "why", "where", "when"}
	ext := filepath.Ext(filename)

	title := strings.Title(
		strings.Replace(
			strings.Replace(
				strings.TrimSuffix(filename, ext),
				"_", " ", -1),
			"-", " ", -1),
	)

	firstWord := strings.Fields(title)[0]
	if containsWord(questionWords, strings.ToLower(firstWord)) {
		title += "?"
	} else {
		title += "!"
	}

	return title
}

func containsWord(words []string, candidate string) bool {
	for _, word := range words {
		if word == candidate {
			return true
		}
	}
	return false
}

func pushButton(buttonId string) (string, error) {
	if !containsWord(AvailableButtonIds(), buttonId) {
		return "", fmt.Errorf("Could not find button with id: %q", buttonId)
	}

	uuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	cmdPath := path.Join(buttonsDir, buttonId)
	fmt.Printf("Running script %v\n", cmdPath)
	cmd := exec.Command(cmdPath)

	// open the out file for writing
	outfile, err := os.Create(path.Join(logsDir, fmt.Sprintf("%v-%v-%v.log", time.Now().UTC().Unix(), uuid, buttonId)))
	if err != nil {
		return "", err
	}
	cmd.Stdout = outfile

	go func() {
		defer outfile.Close()

		err = cmd.Start()
		if err != nil {
			return
		}
		cmd.Wait()
	}()

	return uuid.String(), nil
}

func AvailableButtonIds() []string {
	buttons := AvailableButtons()
	ids := make([]string, len(buttons))
	for _, button := range buttons {
		ids = append(ids, button.Id)
	}
	return ids
}

func loadParameters(filename string) ([]ParameterDef, error) {
	parameters := []ParameterDef{}

	parametersFile := path.Join(buttonsDir, filename+parametersSuffix)
	if !FileExists(parametersFile) {
		return parameters, nil
	}
	bytes, err := ioutil.ReadFile(parametersFile)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(bytes), "\n") {
		if strings.TrimSpace(line) == "" {
			continue // Last line can be empty, I don't care.
		}
		parameter := ParameterDef{}
		components := strings.Split(line, ",")
		parameter.Name = strings.TrimSpace(components[0])
		if parameter.Name == "" {
			return nil, fmt.Errorf("%q is not a valid parameter name for %v", parameter.Name, filename)
		}
		parameter.Type = parameterTypeString
		if len(components) > 1 {
			parameter.Type = strings.TrimSpace(components[1])
		}
		if !containsWord(knownParameterTypes, parameter.Type) {
			return nil, fmt.Errorf("%q is not a valid parameter type for %v", parameter.Name, filename)
		}

		parameter.Description = ""
		if len(components) > 2 {
			parameter.Description = strings.TrimSpace(components[2])
		}

		parameter.Details = make(map[string]interface{})
		if parameter.Type == parameterTypeChoice {
			choices, err := loadChoices(filename, parameter.Name)
			if err != nil {
				return nil, fmt.Errorf("Could not load choices for %v", err)
			}
			parameter.Details["choices"] = choices
		}

		parameters = append(parameters, parameter)
	}

	return parameters, nil

}

func loadChoices(filename string, parameterName string) ([]string, error) {
	return []string{"hello", "yes"}, nil
}
