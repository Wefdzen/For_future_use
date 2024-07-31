package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv" // тут не используется
)

func main() {
	fmt.Println(cfg.PGaddress)
	urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", cfg.PGuser, cfg.PGpassword, cfg.PGaddress, cfg.PGPort, cfg.PGdbname)
	conn, err := pgx.Connect(context.Background(), urlToDataBase)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var (
		id    int
		title string
		price int
	)
	err = conn.QueryRow(context.Background(), "SELECT id, title, price FROM book WHERE id=$1", 1).Scan(&id, &title, &price) //ctx затычка
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(id, title, price)
}

// init читает setting.cfg
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
