package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
  "io/ioutil"
  "bytes"
  "sync"
	"os"
)

/* API.ai stuff */
type WebhookResult struct {
	Action string
  Parameters WebhookParameters 
  Fulfillment WebhookFulfillment
  ResolvedQuery string 
}

type WebhookParameters struct {
  City string
  DateTime string
  GivenName string `json:"given-name"`
}

type WebhookMeta struct {
  IntentName string
}

type WebhookReq struct {
	Result WebhookResult
  Metadata WebhookMeta
  Id string
}

type WebhookFulfillment struct {
	Speech string
}

type WebhookRes struct {
	Fulfillment WebhookFulfillment
}


//#TODO Implement auth tokens 
//*************** AUTH *****************//
func authorization(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Println("Executing middlewareOne")
    next.ServeHTTP(w, r)
  })
}


func textPost(w http.ResponseWriter, r *http.Request) {
 log.Println("Executing finalHandler")
  w.Write([]byte("OK"))
}


func generateToken(w http.ResponseWriter, r *http.Request) {

  w.Header().Set("Content-Type", "application/json")

  decoder := json.NewDecoder(r.Body)

  var t NewTokenStruct
  err := decoder.Decode(&t)

  if err != nil {
    fmt.Println(err)
    //return err 
    fmt.Fprintf(w, "{ \"success\": false, \"err\": { \"message\": %s, \"code\": 10 }", "Invalid Post Parameters")

  } else {

    var key interface{}
    key = []byte("EB32ODSKJN234KJNDSKJSODF89N")
    token, err := signToken(t, key)

    if err != nil {
      fmt.Println(err)
      //return fatal err
      fmt.Fprintf(w, "{ \"success\": false, \"err\": { \"message\": %s, \"code\": 10 } }", "Invalid Post Parameters")
    } else {
      fmt.Println(token)
      //return token 
      fmt.Fprintf(w, "{ \"success\": true, \"token\": \"%s\" }" , token)
    }
  }
}
//*************** END AUTH *****************//

//handles making a simple request 
func simple_req(method string, url string, headers map[string]string, body map[string]string, response interface{}) error {

  tr := &http.Transport { MaxIdleConnsPerHost: 10 }

  client := &http.Client{ Transport: tr }

  var jsonString []byte
  if body != nil {
    var err error
    jsonString, err = json.Marshal(body)
    if err != nil {

    }

  } else { jsonString = []byte("") }

  r, err := http.NewRequest(method, url, bytes.NewBuffer(jsonString))

  for k, v := range headers {
    r.Header.Add(k, v)
  }

  if err != nil {
    return err
  }

  res, err := client.Do(r)

  if err != nil {
    return err
  }

  defer res.Body.Close()
  //b, _ := ioutil.ReadAll(res.Body)
  //fmt.Println(string(b))
  return json.NewDecoder(res.Body).Decode(response)
}

func webhook_handler(rw http.ResponseWriter, request* http.Request) {

	rw.Header().Set("Content-Type", "application/json")

  decoder := json.NewDecoder(request.Body)

	var t WebhookReq
	err := decoder.Decode(&t)

	if err != nil {
		fmt.Println(err)
	}

  fmt.Println(t)

  var response string
  // reads action from api.ai bot and sends back correct api response
  switch t.Result.Action {
    case "reddit_top_post":
      response = "The top post on reddit right now is: " + getRedditTopPost()
    case "current_weather":
      response = getCurrentWeather(t.Result.Parameters.City)
    case "current_bbc_news":
      var res []NewsArticles = getNews("bbc-news")
      if res == nil {
        response = "I could not get news from the BBC right now"
      } else {
        response = "From the BBC news, " + res[0].Title + ". " + res[0].Description
      }
    case "current_cnn_news":
      var res []NewsArticles = getNews("cnn")
      if res == nil {
        response = "I could not get news from CNN right now"
      } else {
        response = "From CNN, " + res[0].Title + ". " + res[0].Description
      }
    case "current_tech_news":
      var res []NewsArticles = getNews("techcrunch")
      if res == nil {
        response = "I could not get news techcrunch right now"
      } else {
        response = "From TechCrunch, " + res[0].Title + ". " + res[0].Description 
      }
    case "joke":
      response = getJoke()
    case "say_my_name":
      var u User
      getUser(&u)
      if u.Name == "" {
         response = "I do not know your name yet"
      } else {
         response = "Your name is " + u.Name
      }
    case "save_my_name":
      var u User
      getUser(&u)
      u.Name = t.Result.Parameters.GivenName
      fmt.Println(u.Name)
      go setUser(u)
      response = "Hi " + t.Result.Parameters.GivenName + " nice to meet you"
    case "current_time":
      //get the current time #TODO
    case "set_reminder":
      //add a reminder 
      response = queueTwilio(t.Result.Parameters.DateTime, t.Result.ResolvedQuery)
    default:
      response = "Sorry, I did not get what you said."
  }

	fmt.Fprintf(rw, "{ \"speech\": \"%s\" }", response)
}

// ******************* Fetching user data ****************//
type User struct {
  Name string `json:"name"`
}

// thread safe get/set functions for user data file
var mutex = &sync.Mutex{}

func getUser(u *User) {
	mutex.Lock()
    file, _ := ioutil.ReadFile("./user_data.json")
    mutex.Unlock()
    json.Unmarshal(file, &u)
}

func setUser(u User) {
	mutex.Lock()
    user_json, _ := json.Marshal(u)
    ioutil.WriteFile("./user_data.json", user_json, 0777)
    mutex.Unlock()
}

func init() {
  // read from user_data.json file for state 
    // if user data file does not exist make one
    _, err := os.Stat("./user_data.json")
    if os.IsNotExist(err) {
        var file, _ = os.Create("./user_data.json")
        defer file.Close()
        empty := []byte{'{', '}'}
        ioutil.WriteFile("./user_data.json", empty, 0777)
    }

}
// ******************* End of user data ****************//

func main() {
	// var u User
	// go getUser(&u)
	// u.Name = "John"
	// go setUser(u)

  //textHandler := http.HandlerFunc(textPost)
	//http.Handle("/api/text", authorization(textHandler))
	http.HandleFunc("/api/speech", webhook_handler)

  http.HandleFunc("/generate_api_token", generateToken)

  fmt.Println("Listening on port 8080...\n")

  log.Fatal(http.ListenAndServe(":8080", nil))

	//call when we generate API key 
	/*out, err := signToken() 

	if err != nil {
		fmt.Println("Err")
	} else {
		err := validateToken(out) 
		if err != nil {
			fmt.Println("err")
		} 
	}
	*/
}
