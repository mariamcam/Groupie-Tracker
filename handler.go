package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Album struct {
	Id           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Location     string              `json:"location"`
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

func (u *Album) getRel(n int) {
	t := *dateLoc
	u.Relation = t.Index[n].Datesandlocations
	for k := range u.Relation {
		if strings.Contains(k, "_") {
			temp := strings.Title(strings.Replace(k, "_", " ", -1))
			u.Relation[temp] = u.Relation[k]
			delete(u.Relation, k)
		} else {
			temp := strings.Title(k)
			u.Relation[temp] = u.Relation[k]
			delete(u.Relation, k)
		}
	}
}

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
	// fmt.Println(data.Index[0])
}

func TakeData() (*[]Album, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data []Album
	jsonerr := json.Unmarshal(responseData, &data)
	if jsonerr != nil {
		return nil, err
	}
	for i := range data {
		data[i].getRel(i)
	}
	return &data, nil
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		err := errors.New("404 Page not found")
		ErrorPage(w, err, 404)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		err := errors.New("405" + "\n" + "Method not allowed")
		ErrorPage(w, err, 405)
		return
	}

	if errorOfData != nil || errorOfRelation != nil {
		err := errors.New("500 Internal Server Error")
		ErrorPage(w, err, 500)
		return
	}

	ts, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		err := errors.New("500 Internal Server Error")
		ErrorPage(w, err, 500)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		ErrorPage(w, err, 500)
		return
	}
}

func Result(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		err := errors.New("405" + "\n" + "Method not allowed")
		ErrorPage(w, err, 405)
		return
	}
	myid, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		err := errors.New("404 Page not found")
		ErrorPage(w, err, 404)
		return
	}

	if myid > 52 || myid < 1 {
		err := errors.New("404 Page not found")
		ErrorPage(w, err, 404)
		return
	}

	t := *data
	index := myid - 1

	tr, err := template.ParseFiles("./templates/result.html")
	if err != nil {
		err := errors.New("500 Internal Server Error")
		ErrorPage(w, err, 500)
		return
	}

	err = tr.Execute(w, t[index])
	if err != nil {
		ErrorPage(w, err, 500)
		return
	}
}

func ErrorPage(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	t, err1 := template.ParseFiles("./templates/error.html")
	if err1 != nil {
		http.Error(w, "500"+"\n"+"Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, err)
	return
}
