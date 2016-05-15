package pushabutton

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/nu7hatch/gouuid"
)

const (
	parametersSuffix    = ".parameters"
	choicesSuffix       = ".choices"
	parameterTypeString = "string"
	parameterTypeChoice = "choice"
)

var knownParameterTypes = []string{parameterTypeChoice, parameterTypeString}

func StartServerOrCrash(addr string) {
	http.Handle("/static/", http.StripPrefix("/static/", ServeAsset{}))
	http.Handle("/api/push/", http.StripPrefix("/api/push/", PostPush{}))
	http.Handle("/api/buttons", http.StripPrefix("/api/buttons", GetButtons{}))
	http.Handle("/api/logs", http.StripPrefix("/api/logs", GetLogs{}))

	http.Handle("/log/", http.StripPrefix("/log/", ServeLog{}))
	http.Handle("/logs", ServeLogIndex{})
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
	tmpl, err := loadTemplate("index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Cannot display page.")
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	err = tmpl.Execute(w, map[string]string{})
}

type PostPush struct{}

func (handler PostPush) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buttonId := r.URL.Path

	var reqPayload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not decode JSON.")
		fmt.Printf("Could not decode JSON: %v\n", err)
		return
	}
	genericArgs, ok := reqPayload["pushArguments"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Parameter pushArguments missing.")
		fmt.Printf("Parameter pushArguments missing")
		return
	}

	var scriptArguments []string
	for _, arg := range genericArgs.([]interface{}) {
		stringArg, ok := arg.(string)
		if !ok {
			fmt.Fprintf(w, "Arguments need to be strings.")
			fmt.Printf("Arguments need to be strings")
			return
		}
		scriptArguments = append(scriptArguments, stringArg)
	}
	fmt.Println("Got arguments", scriptArguments)

	pushId, err := pushButton(buttonId, scriptArguments)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not push button.")
		fmt.Printf("Could not push button: %v\n", err)
		return
	}

	responseData := map[string]string{"buttonId": buttonId, "pushId": pushId}
	payload, err := json.Marshal(responseData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not encode payload as JSON.")
		fmt.Printf("Could not encode payload as JSON: %v\n", err)
		return
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
		return
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

		if shouldIgnoreFileName(filename) {
			return nil
		}

		if strings.Contains(filename, parametersSuffix) || strings.Contains(filename, choicesSuffix) {
			fmt.Printf("skipping parameter file %v of %v\n", candidate, buttons[len(buttons)-1].Id)
			return nil
		}

		buttonId := filename
		title := formatTitle(filename)
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

func formatTitle(filename string) string {
	questionWords := []string{"how", "what", "who", "why", "where", "when"}
	ext := filepath.Ext(filename)

	title := strings.Title(strings.TrimSuffix(filename, ext))

	for _, sep := range []string{".", "-", "_"} {
		title = strings.Replace(title, sep, " ", -1)
	}

	firstWord := strings.Fields(title)[0]
	if containsWord(questionWords, strings.ToLower(firstWord)) {
		title += "?"
	} else {
		title += "!"
	}

	return strings.TrimSpace(title)
}

func containsWord(words []string, candidate string) bool {
	for _, word := range words {
		if word == candidate {
			return true
		}
	}
	return false
}

func pushButton(buttonId string, scriptArguments []string) (string, error) {
	if !containsWord(AvailableButtonIds(), buttonId) {
		return "", fmt.Errorf("Could not find button with id: %q", buttonId)
	}

	pushUuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	pushId := pushUuid.String()

	scriptCall := append([]string{path.Join(buttonsDir, buttonId)}, scriptArguments...)
	now := time.Now().UTC()

	if err := logButtonPush(buttonId, scriptCall, pushId, now); err != nil {
		return "", fmt.Errorf("Could not log button push: %v", err)
	}

	if err := runScriptForButton(buttonId, scriptCall, pushId, now); err != nil {
		return "", fmt.Errorf("Could not run button script: %v", err)
	}

	return pushId, nil
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
				return nil, fmt.Errorf("Could not load choices for %v: %v", parameter.Name, err)
			}
			parameter.Details["choices"] = choices
		}

		parameters = append(parameters, parameter)
	}

	return parameters, nil
}

func loadChoices(filename string, parameterName string) ([]string, error) {
	var choices []string

	choicesScriptPrefix := path.Join(buttonsDir, filename+parametersSuffix+"."+parameterName+choicesSuffix)
	fmt.Printf("Determining choices for filename %v and parameter %v by running the script with prefix: %v\n", filename, parameterName, choicesScriptPrefix)
	choicesScript := findSingleScriptWithPrefix(choicesScriptPrefix)
	if choicesScript == "" {
		return nil, fmt.Errorf("ERROR: Could not read choices for parameter %v, could not find suitable script", parameterName)
	}
	var stdout bytes.Buffer

	fmt.Printf("Running script %v\n", choicesScript)
	cmd := exec.Command(choicesScript)
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("ERROR: Could not read choices for parameter %v, could not start script: %v", parameterName, err)
	}

	err = cmd.Wait()
	if exitError, ok := err.(*exec.ExitError); ok {
		waitStatus := exitError.Sys().(syscall.WaitStatus)
		return nil, fmt.Errorf("ERROR: Could not read choices for parameter %v, script returned non-zero exit code (%v): %v", parameterName, waitStatus.ExitStatus(), err)
	} else if err != nil {
		return nil, fmt.Errorf("ERROR: Could not read choices for parameter %v, script failed: %v", parameterName, err)
	}

	for _, line := range strings.Split(string(stdout.Bytes()), "\n") {
		choices = append(choices, strings.TrimSpace(line))
	}

	return choices, nil
}

func findSingleScriptWithPrefix(prefix string) string {
	var candidates []string
	filepath.Walk(buttonsDir, func(candidate string, info os.FileInfo, err error) error {
		if candidate == buttonsDir {
			return nil
		}
		if !strings.HasPrefix(candidate, prefix) {
			return nil
		}
		filename := path.Base(candidate)
		if shouldIgnoreFileName(filename) {
			return nil
		}
		if info.Mode()&os.ModePerm&0100 == 0 {
			fmt.Printf("Ignoring non-executable: %v\n", candidate)
			return nil
		}
		candidates = append(candidates, candidate)

		return nil
	})
	if len(candidates) != 1 {
		fmt.Printf("Found %v candidates for prefix, so ignoring them: %v\n", len(candidates), prefix)
		return ""
	}
	return candidates[0]
}

func logButtonPush(buttonId string, scriptCall []string, uuid string, now time.Time) error {
	records := []string{
		strconv.FormatInt(now.Unix(), 10),
		now.Format(time.RFC3339),
		uuid,
		buttonId,
		strings.Join(scriptCall, " "),
	}
	f, err := os.OpenFile(logfilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	csvWriter := csv.NewWriter(f)
	if err = csvWriter.Write(records); err != nil {
		return err
	}
	if err = csvWriter.Error(); err != nil {
		return err
	}

	csvWriter.Flush()
	if err = csvWriter.Error(); err != nil {
		return err
	}

	return nil
}

func runScriptForButton(buttonId string, scriptCall []string, pushId string, now time.Time) error {
	fmt.Printf("Running script %v\n", strings.Join(scriptCall, " "))
	cmd := exec.Command(scriptCall[0], scriptCall[1:]...)

	// open the out file for writing
	pushLogPath := path.Join(logsDir, fmt.Sprintf("%v-%v-%v.log", now.Unix(), pushId, buttonId))
	outfile, err := os.Create(pushLogPath)
	if err != nil {
		return err
	}
	cmd.Stdout = outfile
	cmd.Stderr = outfile

	go func() {
		defer outfile.Close()

		err := cmd.Start()
		if err != nil {
			logPushResult(outfile, fmt.Sprintf("ERROR: Command did not even start: %v", err))
		}
		err = cmd.Wait()
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			logPushResult(outfile, fmt.Sprintf("FAILURE: Exited with exit code (%v): %v", waitStatus.ExitStatus(), exitError))
			return

		} else if err != nil {
			logPushResult(outfile, fmt.Sprintf("ERROR: Script did not run: %v", err))
			return
		}

		logPushResult(outfile, "SUCCESS: Command exited without errors")
	}()

	return nil
}

func logPushResult(outfile io.Writer, statusline string) {
	fmt.Fprintf(outfile, "\n\n============================\n"+statusline+"\n============================\n")
}

type ServeLogIndex struct{}

func (handler ServeLogIndex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := loadTemplate("logs.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Cannot display page.")
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	err = tmpl.Execute(w, map[string]string{})
}

func loadTemplate(filename string) (*template.Template, error) {
	data, err := Asset(path.Join(assetsDir, filename))
	if err != nil {
		return nil, fmt.Errorf("Could not find template: %v", err)
	}
	indexHtml := string(data)
	tmpl, err := template.New("index.html").Parse(indexHtml)
	if err != nil {
		return nil, fmt.Errorf("Could not parse template: %v", err)
	}

	return tmpl, nil
}

type LogSummaryEntry struct {
	PushId      string
	ButtonId    string
	Timestamp   string
	DateTimeUTC string
	Title       string
	Cmd         string
}

type GetLogs struct{}

func (handler GetLogs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	entries, err := AvailableLogs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not read logs.")
		fmt.Printf("Could not read logs: %v\n", err)
		return
	}
	payload, err := json.Marshal(entries)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not encode payload as JSON.")
		fmt.Printf("Could not encode payload as JSON: %v\n", err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(payload)
}

func AvailableLogs() ([]LogSummaryEntry, error) {
	var entries []LogSummaryEntry

	csvFile, err := os.OpenFile(logfilePath, os.O_RDONLY, 0644)
	if os.IsNotExist(err) {
		return entries, nil
	}
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(csvFile)

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if len(record) != 5 {
			fmt.Printf("Ignoring record: %v", record)
			continue
		}
		entries = append(entries, LogSummaryEntry{
			Timestamp:   strings.TrimSpace(record[0]),
			DateTimeUTC: strings.TrimSpace(record[1]),
			PushId:      strings.TrimSpace(record[2]),
			ButtonId:    strings.TrimSpace(record[3]),
			Cmd:         strings.TrimSpace(record[4]),
			Title:       formatTitle(strings.TrimSpace(record[3])),
		})
	}

	return entries, nil
}

type LogEntry struct {
	LogSummaryEntry
	Stdouterr string
}

type ServeLog struct{}

func (handler ServeLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pushId := r.URL.Path
	fmt.Printf("Getting log for %v\n", pushId)

	entry := findOrGenerateLogEntry(pushId)

	tmpl, err := loadTemplate("logEntry.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Cannot display page.")
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	autorefreshInSec, err := strconv.Atoi(r.URL.Query().Get("autorefresh"))
	if err != nil {
		autorefreshInSec = 0
	}
	tmpl.Execute(w, map[string]interface{}{
		"autorefreshInSec": autorefreshInSec,
		"entry":            entry,
	})
}

func findOrGenerateLogEntry(pushId string) LogEntry {
	entry := LogEntry{
		LogSummaryEntry: LogSummaryEntry{
			PushId:   pushId,
			ButtonId: "N/A",
			Cmd:      "N/A",
			Title:    "N/A",
		},
		Stdouterr: "N/A",
	}
	logs, err := AvailableLogs()
	if err != nil {
		fmt.Printf("Cannot read logs: %v\n", err)
		logs = []LogSummaryEntry{}
	}
	for _, summary := range logs {
		if summary.PushId == pushId {
			entry.LogSummaryEntry = summary
			cmdOutputPath := path.Join(
				logsDir,
				fmt.Sprintf("%v-%v-%v.log", summary.Timestamp, pushId, summary.ButtonId),
			)
			stdouterr, err := ioutil.ReadFile(cmdOutputPath)
			if err != nil {
				fmt.Printf("Could not read stdout/stderr of %v: %v\n", pushId, err)
			} else {
				entry.Stdouterr = string(stdouterr)
			}
			break
		}
	}

	return entry
}

func shouldIgnoreFileName(filename string) bool {
	if strings.HasPrefix(filename, ".") || strings.HasSuffix(filename, "~") || strings.HasSuffix(filename, "#") || strings.TrimSpace(filename) == "" {
		fmt.Printf("skipping hidden/temporary/suspicous file %q\n", filename)
		return true
	}

	return false
}
