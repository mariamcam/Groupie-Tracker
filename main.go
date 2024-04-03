package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Artists struct {
	Id           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Locations    string              `json:"locations"`
	ConcertDates string              `json:"concertDates"`
	Relation     map[string][]string `json:"relation"`
}
type Dl struct {
	Index []struct {
		Datesandlocations map[string][]string `json:"datesLocations"`
	}
}

var (
	data, errorOfData        = TakeData()
	dateLoc, errorOfRelation = GetRelation()
)

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("templates/css/"))
	mux.Handle("/css/", http.StripPrefix("/css/", fs))
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/result/", Result)
	fmt.Println("serveur ouvert sur Go... http://localhost:3534.")
	fmt.Println(http.ListenAndServe(":3534", mux))

}

// GetRelation récupère les données de relation depuis l'API externe.
func GetRelation() (*Dl, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data Dl
	jsonerr := json.Unmarshal(responseData, &data)
	if jsonerr != nil {
		return nil, err
	}

	return &data, nil
}

// TakeData récupère les données d'artistes depuis l'API externe.
func TakeData() (*[]Artists, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data []Artists
	jsonerr := json.Unmarshal(responseData, &data)
	if jsonerr != nil {
		return nil, err
	}

	for i := range data {
		data[i].Relation = dateLoc.Index[i].Datesandlocations
	}
	return &data, nil
}

// Home est le gestionnaire pour la page d'accueil.
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		err := errors.New("404 Page non trouvée")
		ErrorPage(w, err, 404)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		err := errors.New("405" + "\n" + "Méthode non autorisée")
		ErrorPage(w, err, 405)
		return
	}

	if errorOfData != nil || errorOfRelation != nil {
		err := errors.New("500 Erreur interne du serveur")
		ErrorPage(w, err, 500)
		return
	}

	ts, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		err := errors.New("500 Erreur interne du serveur")
		ErrorPage(w, err, 500)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		ErrorPage(w, err, 500)
		return
	}
}

// Result est le gestionnaire pour la page de résultat.
func Result(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		err := errors.New("405" + "\n" + "Méthode non autorisée")
		ErrorPage(w, err, 405)
		return
	}
	myid, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		err := errors.New("404 Page non trouvée")
		ErrorPage(w, err, 404)
		return
	}

	if myid > 52 || myid < 1 {
		err := errors.New("404 Page non trouvée")
		ErrorPage(w, err, 404)
		return
	}

	t := *data
	index := myid - 1

	tr, err := template.ParseFiles("./templates/result.html")
	if err != nil {
		err := errors.New("500 Erreur interne du serveur")
		ErrorPage(w, err, 500)
		return
	}

	err = tr.Execute(w, t[index])
	if err != nil {
		ErrorPage(w, err, 500)
		return
	}
}

// ErrorPage affiche une page d'erreur en fonction du code d'erreur donné.
func ErrorPage(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	t, err1 := template.ParseFiles("./templates/error.html")
	if err1 != nil {
		http.Error(w, "500"+"\n"+"Erreur interne du serveur", http.StatusInternalServerError)
		return
	}
	t.Execute(w, err)
	// Supprimer la ligne panic(err)
}
