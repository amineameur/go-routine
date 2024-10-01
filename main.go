package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const apiKey = "Enter_your_key_here"

var cities = []string{"London", "Paris", "Tunisia", "Tokyo", "Berlin"}

type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func fetchWeather(city string) (WeatherResponse, error) {
	var weatherData WeatherResponse
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return weatherData, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &weatherData)
	return weatherData, err
}

func fetchWeatherSequential() {
	fmt.Println("Fetching weather data sequentially:")
	start := time.Now()

	for _, city := range cities {
		weather, err := fetchWeather(city)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Printf("City: %s, Temp: %.2f°C\n", weather.Name, weather.Main.Temp)
	}

	fmt.Printf("Time taken without Go routine: %v\n\n", time.Since(start))
}

func fetchWeatherConcurrent() {
	fmt.Println("Fetching weather data with Go routines:")
	start := time.Now()
	var wg sync.WaitGroup

	for _, city := range cities {
		wg.Add(1)
		go func(city string) {
			defer wg.Done()
			weather, err := fetchWeather(city)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("City: %s, Temp: %.2f°C\n", weather.Name, weather.Main.Temp)
		}(city)
	}

	wg.Wait()
	fmt.Printf("Time taken with Go routine: %v\n", time.Since(start))
}

func main() {

	fetchWeatherConcurrent()
}
