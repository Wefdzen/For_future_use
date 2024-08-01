package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

// struct User
type User struct {
	Name  string
	Age   uint16
	Money int64
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func account(w http.ResponseWriter, r *http.Request) { // w это типо страничка request запросы с странички!!!
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		password := r.FormValue("password")

		urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", cfg.PGuser, cfg.PGpassword, cfg.PGaddress, cfg.PGPort, cfg.PGdbname)
		conn, err := pgx.Connect(context.Background(), urlToDataBase)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		conn.Exec(context.Background(), "INSERT INTO appli (name, password) VALUES ($1, $2)", name, password)
		// file, err := os.OpenFile("output.txt", os.O_WRONLY|os.O_APPEND, 0600)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// defer file.Close()
		// file.WriteString(fmt.Sprintf("Name: %s, password: %s\n", name, password))

	}

	tmpl, err := template.ParseFiles("templates/account_page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Dev := User{"Dima", 19, 200}
	err = tmpl.Execute(w, Dev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func main() {

	http.HandleFunc("/", mainPage)
	http.HandleFunc("/account/", account)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func init() {
	file, err := os.Open("config.cfg")
	if err != nil {
		fmt.Println("Error open .cfg", err)
		panic("Can't open the file \"setting.cfg\"")
	}
	defer file.Close()

	fileInfo, _ := file.Stat()                   // получаю стату файла для его размера
	readSetting := make([]byte, fileInfo.Size()) // делаю такого же размера переменную
	_, err = file.Read(readSetting)
	if err != nil {
		panic("can't read file")
	}
	// fmt.Println(string(readSetting))  работает

	err = json.Unmarshal(readSetting, &cfg) //unmarshal и json в обьект marshal из object in json
	if err != nil {
		panic("json err")
	}
}

type setting struct { // должен повторять структуру json
	PGaddress  string
	PGpassword string
	PGuser     string
	PGdbname   string
	PGPort     string
}

var (
	cfg setting
)
