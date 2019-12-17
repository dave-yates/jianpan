package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/dave-yates/jianpan/chinese"
	"github.com/dave-yates/jianpan/db"
	"github.com/dave-yates/jianpan/importer"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client mongo.Client

func main() {

	db.SetupConfig()

	//context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newClient, err := mongo.NewClient(options.Client().ApplyURI(db.Config.URI))
	client = *newClient
	defer client.Disconnect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	setupDatabase(ctx)

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

func setupDatabase(ctx context.Context) {

	//setup mongodb
	err := db.InitDB(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	err = importer.Import(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func getTranslation(w http.ResponseWriter, r *http.Request) {

	input := r.URL.Query().Get("input")

	//context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	output, _ := chinese.Translate(ctx, input)

	fmt.Println(string(output))

	w.Write(output)
}
