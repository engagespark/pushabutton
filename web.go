package pushabutton

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

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

	responseData := map[string]string{"buttonId": buttonId}
	payload, err := json.Marshal(responseData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not encode payload as JSON.")
		fmt.Printf("Could not encode payload as JSON: %v\n", err)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(payload)
	fmt.Printf("Pressed " + buttonId + "!\n")
}

type Button struct {
	Id    string
	Title string
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
		} else {
			fmt.Println(info.Mode())
		}

		filename := path.Base(candidate)

		if strings.TrimSpace(filename) == "" {
			fmt.Printf("skipping suspicous whitespace file %q\n", candidate)
			return nil
		}

		buttonId := filename
		title := generateTitle(filename)

		buttons = append(buttons, Button{
			Id:    html.EscapeString(buttonId),
			Title: html.EscapeString(title),
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
