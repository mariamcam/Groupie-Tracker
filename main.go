package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()                          // Хранит в себе все URL-ы и зависимости приложения
	fs := http.FileServer(http.Dir("templates/css/"))  // handle для нахождения css файлы
	mux.Handle("/css/", http.StripPrefix("/css/", fs)) // удаление префикса "templates" для доступа именно к css либо html файлам
	mux.HandleFunc("/", Home)                          // handle домашней странички
	mux.HandleFunc("/result/", Result)
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	fmt.Println(http.ListenAndServe(":8080", mux))

	// Data := "https://groupietrackers.herokuapp.com/api/artists"
	// resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	// if err != nil {
	// 	return
	// }
	// jsonerr := json.Unmarshal(resp, &album)

	// fmt.Print(data[0].getIdNameImage)
}

// func (a *Album) getIdNameImage() string {
// 	return "asdasd"
// }
