package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

type Status struct {
	Status Weather `json:"status"`
}

type Weather struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Hasil struct {
	Water       int `json:"water"`
	Wind        int `json:"wind"`
	StatusWater string
	StatusWind  string
}

func main() {
	address := "localhost:9090"
	http.HandleFunc("/", mainPage)
	http.Handle("/asset/", http.StripPrefix("/asset/", http.FileServer(http.Dir("asset"))))
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("Error running service: ", err)
	}

}

func mainPage(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("weather.json")
	if err != nil {
		fmt.Print(err)
	}

	var obj Status

	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}

	writeToFile()

	var hasil Hasil
	hasil.Water = obj.Status.Wind
	hasil.Wind = obj.Status.Wind

	if obj.Status.Wind <= 6 {
		hasil.StatusWind = "Aman"
	}
	if obj.Status.Wind >= 7 && obj.Status.Wind <= 15 {
		hasil.StatusWind = "Siaga"
	}
	if obj.Status.Wind > 15 {
		hasil.StatusWind = "Bahaya"
	}

	if obj.Status.Water <= 5 {
		hasil.StatusWater = "Aman"
	}
	if obj.Status.Water >= 6 && obj.Status.Water <= 8 {
		hasil.StatusWater = "Siaga"
	}
	if obj.Status.Water > 8 {
		hasil.StatusWater = "Bahaya"
	}

	t, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	err = t.ExecuteTemplate(w, "index.html", hasil)
	if err != nil {
		panic(err)
	}
}

func writeToFile() {
	wind := rand.Intn(100)
	water := rand.Intn(100)
	dataStatus := Status{
		Status: Weather{Wind: wind,
			Water: water},
	}
	file, _ := json.MarshalIndent(dataStatus, "", " ")

	_ = ioutil.WriteFile("weather.json", file, 0644)
}
