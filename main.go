package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Clima struct {
	Localizacao struct {
		Nome string `json:"name"`
		Pais string `json:"country"`
	} `json:"location"`
	Atual struct {
		TempC    float64 `json:"temp_c"`
		Condicao struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Previsao struct {
		Previsaodia []struct {
			Hora []struct {
				Tempo    int64   `json:"time_epoch"`
				TempC    float64 `json:"temp_c"`
				Condicao struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceDeChuva float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	q := "Chapeco"

	if len(os.Args) > 2 {
		q = os.Args[1]
	}
	resposta, erro := http.Get("https://api.weatherapi.com/v1/forecast.json?key=a1f040c66f1149e98f0181223240108&q=" + q + "&days1&aqi=no&alerts=no")
	if erro != nil {
		log.Fatal(erro)
	}
	defer resposta.Body.Close()

	if resposta.StatusCode != 200 {
		log.Fatal("API não está respondendo")
	}

	body, erro := io.ReadAll(resposta.Body)
	if erro != nil {
		log.Fatal(erro)
	}

	var clima Clima
	erro = json.Unmarshal(body, &clima)
	if erro != nil {
		log.Fatal(erro)
	}

	local, atual, horas := clima.Localizacao, clima.Atual, clima.Previsao.Previsaodia[0].Hora

	fmt.Printf("%s, %s: %.0fC, %s\n",
		local.Nome,
		local.Pais,
		atual.TempC,
		clima.Atual.Condicao.Text,
	)

	for _, hora := range horas {
		date := time.Unix(hora.Tempo, 0)

		if date.Before(time.Now()) {
			continue
		}
		mensagem := fmt.Sprintf(
			"%s - %.0fC, %0.f%%, %s\n",
			date.Format("15:04"),
			hora.TempC,
			hora.ChanceDeChuva,
			hora.Condicao.Text,
		)

		if hora.ChanceDeChuva < 40 {
			fmt.Print(mensagem)
		} else {
			color.Red(mensagem)
		}
	}
}
