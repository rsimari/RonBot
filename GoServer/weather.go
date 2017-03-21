package main 


import (
  "strings"
  "fmt"
)

//****************** START weather ******************//

// Weather structs
type WeatherRes struct {
  Main WeatherMain
  Weather[] WeatherData
}

type WeatherData struct {
  Main string 
  Id int 
  Description string 
}

type WeatherMain struct {
  Temp float32 
  MinTemp float32 `json:"temp_min"`
  MaxTemp float32 `json:"temp_max"`
}

func getCurrentWeather(city string) string {

  response := WeatherRes{}

	var replacer = strings.NewReplacer(" ", "+")
	var urlCity = replacer.Replace(city)
  simple_req("GET", "http://api.openweathermap.org/data/2.5/weather?q=" + urlCity + "&APPID=4e7036fa40c4ae2705533033fa77b0a1", nil, nil, &response)

  if response.Main.Temp > 0 && response.Main.MinTemp > 0 && response.Main.MaxTemp > 0 {

    //Calculate appropriate temps 
  	var fTemp = 1.8*(response.Main.Temp - 273) + 32
  	var fLowTemp = 1.8*(response.Main.MinTemp - 273) + 32
  	var fHighTemp = 1.8*(response.Main.MaxTemp-273) + 32

    //Format temp strings 
  	s := fmt.Sprintf("%.0f", fTemp)
  	sLow := fmt.Sprintf("%.0f", fLowTemp)
  	sHigh := fmt.Sprintf("%.0f", fHighTemp)

    //Parse weather description 
  	var id = response.Weather[0].Id
  	var idGroup = id/100 
  	var description string

  	switch idGroup {
  		case 8:
    		if id % 100 > 0 {
    			if id == 801 || id == 802 {
    				description = "slightly cloudy"
    			} else {
    				description = "cloudy"
     			}
    	   } else {
    			description = "clear"
    	 } 
  		case 9: 
  		  description = "extreme"
  		case 6:
  		  description = "snowy"
  		case 5: 
  		  description = "rainy"
  		case 3:
  		  description = "drizzling"
  		case 2:
  		  description = "thunderstorming"
  		default: 
  		  description = "sunny"
  	}

    //Return weather string 
    return "It is " + description + " in " + city + " right now. The temperature is " + s + " degrees, with a high of " + sHigh + " and a low of " + sLow + "."
    } else { return "I am sorry. I am unable to retrieve the weather." }
}
//****************** END weather ******************//