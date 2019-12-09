package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/dave-yates/jianpan/dictionary"
	"github.com/dave-yates/jianpan/importer"
)

var dict dictionary.Dictionary

func main() {

	dict = importer.ImportData()

	http.HandleFunc("/", handler)
	http.HandleFunc("/keyboard", keyboardHandler)
	http.HandleFunc("/help", helpHandler)
	http.HandleFunc("/translations", getTranslation)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world,  你好世界"))
}

func keyboardHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/keyboard.html")
	t.Execute(w, nil)
}

func helpHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/help.html")
	t.Execute(w, nil)
}

func getTranslation(w http.ResponseWriter, r *http.Request) {

	input := r.URL.Query().Get("input")

	output, _ := dict.Translate(input)

	fmt.Println(string(output))

	w.Write(output)
}
