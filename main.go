package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	// "strconv"
	"text/template"
)

// "github.com/emanoelxavier/openid2go/openid"

//GamesOwned es la estructura de datos a mostrar
type GamesOwned struct {
	Response struct {
		GameCount int `json:"game_count"`
		Games     []struct {
			Appid                  int `json:"appid"`
			PlaytimeForever        int `json:"playtime_forever"`
			PlaytimeWindowsForever int `json:"playtime_windows_forever"`
			PlaytimeMacForever     int `json:"playtime_mac_forever"`
			PlaytimeLinuxForever   int `json:"playtime_linux_forever"`
			Playtime2Weeks         int `json:"playtime_2weeks,omitempty"`
		} `json:"games"`
	} `json:"response"`
}

//RandomGame es una estructura usada para obtener los desafios
type RandomGame struct {
	Appid    int `json:"appid"`
	Desafios string
}

//IDsteam64 es una estructura usada para guardar la id el usuario
type IDsteam64 struct {
	Id64 string
}

func main() {

	http.HandleFunc("/home", ObtenerJuegos)
	http.HandleFunc("/Desafio", arrayJuegos)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))

}

//ObtenerJuegos devuelve los juegos que posee el usuario
func ObtenerJuegos(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("FirstScreen.html")
	if err != nil {
		fmt.Println("Index Template Parse Error: ", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Index Template Execution Error: ", err)
	}

	steamId := r.FormValue("steamID64")

	url := ("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=BFFE48E8F58A9E59161FBBC9D8DD5A2B&steamid=" + steamId + "&format=json")
	client := &http.Client{}
	fmt.Print(url)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Printf("Falló la creación del request a la URL '%s', dando el error %v", url, err.Error())
		os.Exit(1)
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Falló el acceso a la URL '%s', dando el error %v", url, err.Error())
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Falló el acceso al body de la respuesta de '%s', dando el error %v", url, err.Error())
		os.Exit(1)
	}

	error := ioutil.WriteFile("JuegosAdquiridos", body, 0644)
	if error != nil {
		fmt.Println("no se puede escribir el archivo")
	}

}

func arrayJuegos(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))

	archivo := "JuegosAdquiridos"
	var juegos GamesOwned
	var desafiosArray = [...]string{"platina el juego", "Terminalo antes de 3 dias", "Pasalo sin Morir, si es imposible morir, pasalo con una mano."}

	ArrayJuegos, err := ioutil.ReadFile(archivo)

	if err != nil {
		fmt.Println("Archivo inválido en la ruta: '", "' ", err.Error())
		os.Exit(1)
	}

	_ = json.Unmarshal(ArrayJuegos, &juegos)

	randomN := rand.Int() % len(juegos.Response.Games)
	randomNdesafio := rand.Int() % len(desafiosArray)

	var juegoRandom int
	juegoRandom = juegos.Response.Games[randomN].Appid
	desafioAleatorio := desafiosArray[randomNdesafio]

	resultadoFinal := RandomGame{
		Appid:    juegoRandom,
		Desafios: desafioAleatorio,
	}
	tmpl.Execute(w, resultadoFinal)

}
