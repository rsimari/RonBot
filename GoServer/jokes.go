package main 
import (
  "fmt"
)

//****************** START jokes ******************//

type JokeRes struct {
  Type string
  Value JokeData
}

type JokeData struct {
  Joke string
}

func getJoke() string {
  fmt.Println("Getting a joke...")

  response := JokeRes{}

  simple_req("GET", "http://api.icndb.com/jokes/random?limitTo=[nerdy]&firstName=Peter&lastName=Bui&escape=javascript", nil, nil, &response)

  if response.Type != "success" {
    return "Sorry I could not get a joke"
  }

  return response.Value.Joke

}

//***************** END jokes ********************//