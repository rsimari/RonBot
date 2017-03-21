package main

import (
	"fmt"
	"net/http"
	"net/url"
	"encoding/json"
	"log"
  "io/ioutil"
  "bytes"
  "sync"
	"strings"
	"os"

  "time"

	//"net/http/httputil"
)


//************** IRON.IO STUFF *******************//

type Task struct {
        CodeName string `json:"code_name"`
        Payload string `json:"payload"`
        StartAt string `json:"start_at"`
}

type ReqData struct {
        Schedules []*Task `json:"schedules"`
}

func queueTwilio(execTime string, msg string) string {

        if len(execTime) == 8 {
          //add todays calendar date etc.
	loc, _ := time.LoadLocation("America/New_York") //get from file  
          date := time.Now().In(loc)
	dateString := date.Format("2006-01-02")   
	fmt.Println(dateString)
          execTime = string(dateString) + "T" + execTime + "-04:00"
	fmt.Println(execTime)
	newTime, _ := time.Parse("2006-01-02T15:04:05-07:00", execTime)
	newUTCTime := newTime.UTC()

        fmt.Println(newUTCTime.Format("2006-01-02T15:04:05Z"))
	} else {

	//loc, _ := time.LoadLocation("America/New_York")
	execTime := execTime + "-04:00"
	newTime, _ := time.Parse("2006-01-02T15:04:05Z-07:00", execTime)

	newUTCTime := newTime.UTC()

	fmt.Println(newUTCTime.Format("2006-01-02T15:04:05Z"))	

	}

        const token = "IpsW3ZNlFsLl42T2vSqw"
        const project = "58cf2ffd0e2c7300061ad812"

        // Insert our project ID and token into the API endpoint
        target := fmt.Sprintf("http://worker-us-east.iron.io/2/projects/%s/schedules?oauth=%s", project, token)

        // Build the payload
        // The payload is a string to pass information into your worker as part of a task
        // It generally is a JSON-serialized string (which is what we're doing here) that can be deserialized in the worker
        payload := map[string]interface{} {
                "Message" : msg,
                "Phone" : "+14127601315",
        }

        payload_bytes, err := json.Marshal(payload)
        if err != nil {
                panic(err.Error())
        }
        payload_str := string(payload_bytes)

        // Build the task
        task := &Task {
                CodeName: "send_twilio",
                Payload: payload_str,
                StartAt: execTime,
        }

        // Build a request containing the task
        json_data := &ReqData {
                Schedules: []*Task { task },
        }

        json_bytes, err := json.Marshal(json_data)

        if err != nil {
                panic(err.Error())
        }

        json_str := string(json_bytes)

        // Post expects a Reader
        json_buf := bytes.NewBufferString(json_str)

        // Make the request
        resp, err := http.Post(target, "application/json", json_buf)
        if err != nil {
                panic(err.Error())
        }
        defer resp.Body.Close()

        // Read the response
        resp_body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                panic(err.Error())
        }

        // Print the response to STDOUT
        fmt.Println(string(resp_body))

        return "Okay. I will send you a text message to remind you!"
}

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


func spotify_auth() {
  var BASE_URI string = "https://accounts.spotify.com/api/token"
  req, err := http.NewRequest("POST", BASE_URI, nil)
  if err != nil {
    fmt.Println(err)
    return
  }
  req.Header.Set("Authorization", "Basic")
  q := req.URL.Query()
  q.Add("grant_type", "client_credentials")
  q.Add("response_type", "code")
  q.Add("redirect_uri", spotify_redirect)

  req.URL.RawQuery = q.Encode()
  fmt.Println(req.URL.String())
  client := &http.Client{}
  resp, er := client.Do(req)
  if er != nil {
    fmt.Println(er)
  }
  fmt.Println(resp.Body)
}


//************* HACKER NEWS SKILL *******************///

type HackerNewsJSON struct {
  Title string
}

func getHackerNewsTopPost() string {
  var BASE_URI string = "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
  res, err := http.Get(BASE_URI)
  if err != nil {
    fmt.Println(err)
  }
  defer res.Body.Close()

  var resp string = ""
  if res.StatusCode == 200 {
    bodyBytes, _ := ioutil.ReadAll(res.Body)
    var bodyString string = string(bodyBytes)
    var topID string = ""
    var i int = 2
    for bodyString[i] != 44 {
      topID = topID + string(bodyString[i])
      i = i + 1
    }
    response := HackerNewsJSON{}
    simple_req("GET", "https://hacker-news.firebaseio.com/v0/item/" + topID + ".json?print=pretty", nil, nil, &response)
    resp = string(response.Title)
  }
  return resp
}

func getHackerNewsTopJob() string {
  var BASE_URI string = "https://hacker-news.firebaseio.com/v0/jobstories.json"
  res, err := http.Get(BASE_URI)
  if err != nil {
    fmt.Println(err)
  }
  defer res.Body.Close()

  var resp string = ""
  if res.StatusCode == 200 {
    bodyBytes, _ := ioutil.ReadAll(res.Body)
    var bodyString string = string(bodyBytes)
    var topID string = ""
    var i int = 1
    for bodyString[i] != 44 {
      topID = topID + string(bodyString[i])
      i = i + 1
    }
    response := HackerNewsJSON{}
    simple_req("GET", "https://hacker-news.firebaseio.com/v0/item/" + topID + ".json?print=pretty", nil, nil, &response)
    resp = string(response.Title)
  }
  return resp

}


//****************** END HACKER NEWS **********************//

//****************** REDDIT SKILL *************************//

// specific information for post
type RedditChildData struct {
  Title string
}
type RedditPosts struct {
  Data RedditChildData
}
// data holding all posts on page
type RedditPostData struct {
  Children []RedditPosts
  Modhash string
}
// overall json from page
type RedditRes struct {
  Data RedditPostData
  Kind string
}
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

func getRedditTopPost() string {
  response := RedditRes{}
  simple_req("GET", "https://www.reddit.com/r/all.json", nil, nil, &response)
  return response.Data.Children[0].Data.Title
}


//****************** END reddit **********************//


//****************** START news *********************//

type NewsRes struct {
  Source string `json:"source"`
  Status string `json:"status"`
  Articles []NewsArticles `json:"articles"`
}

type NewsArticles struct {
  Author string
  Title string
  Description string
}

func getNews(src string) []NewsArticles {
  fmt.Println("Getting news from " + src + "...");

  response := NewsRes{}

  header := make(map[string]string)
  header["x-api-key"] = "fcad7e1888274b8486f80b4e7435692e"

  simple_req("GET", "https://newsapi.org/v1/articles?source=" + src, header, nil, &response)

  if response.Status != "ok" {
    return nil
  }

  return response.Articles
}

//****************** END news **********************//

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

//****************** START weather ******************//
func getCurrentWeather(city string) string {
  fmt.Println("Getting weather...");

  response := WeatherRes{}

	var replacer = strings.NewReplacer(" ", "+")
	var urlCity = replacer.Replace(city)
  simple_req("GET", "http://api.openweathermap.org/data/2.5/weather?q=" + urlCity + "&APPID=4e7036fa40c4ae2705533033fa77b0a1", nil, nil, &response)


	var fTemp = 1.8*(response.Main.Temp - 273) + 32
	var fLowTemp = 1.8*(response.Main.MinTemp - 273) + 32
	var fHighTemp = 1.8*(response.Main.MaxTemp-273) + 32
	s := fmt.Sprintf("%.0f", fTemp)
	sLow := fmt.Sprintf("%.0f", fLowTemp)
	sHigh := fmt.Sprintf("%.0f", fHighTemp)

	var id = response.Weather[0].Id
	var idGroup = id/100 
	var description string
	fmt.Println(id)	
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

  return "It is " + description + " in " + city + " right now. The temperature is " + s + " degrees, with a high of " + sHigh + " and a low of " + sLow + "."
}

func setTwilioReminder() {

	accountSid := "AC6fe206e4c09c5f10f3f6907d2d4a1498"
	authToken := "b917e2302a83ba7cc98731987628d8c9"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	v := url.Values{}
  	v.Set("To","+14127601315")
  	v.Set("From","+14123608598")
  	v.Set("Body","Brooklyn's in the house!")
  	rb := *strings.NewReader(v.Encode())

	client := &http.Client{}
 
  	req, _ := http.NewRequest("POST", urlStr, &rb)
  	req.SetBasicAuth(accountSid, authToken)
  	req.Header.Add("Accept", "application/json")
  	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

  	// Make request
  	resp, _ := client.Do(req)
	  fmt.Println(resp.Status)


}


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

var spotify_client_id string = "a8c0b2ec2d4542298259a9c6d85dba83"
var spotify_client_secret string = "a01877d1e09245e3a1f22f04b8a9fc1e"
var spotify_redirect string = "https://35.166.199.67:8080/"

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
    case "hacker_news_top_post":
      response = "The top post of hacker news is: " + getHackerNewsTopPost()
    case "hacker_news_top_job":
      response = "The top job post on hacker news is currently: " + getHackerNewsTopJob()
    case "reddit_top_post":
      response = "The top post on reddit right now is: " + getRedditTopPost()
    case "current_weather":
      response = getCurrentWeather(t.Result.Parameters.City)
    case "current_bbc_news":
      var res []NewsArticles = getNews("bbc-news")
      if res == nil {
        response = "I could not get news from the BBC right now"
      } else {
        response = "From the BBC, " + res[0].Title
      }
    case "current_cnn_news":
      var res []NewsArticles = getNews("cnn")
      if res == nil {
        response = "I could not get news from CNN right now"
      } else {
        response = "From CNN, " + res[0].Title
      }
    case "current_tech_news":
      var res []NewsArticles = getNews("techcrunch")
      if res == nil {
        response = "I could not get news techcrunch right now"
      } else {
        response = "From TechCrunch, " + res[0].Title
      }
    case "joke":
      response = getJoke()
    case "current_time":
      //get the current time 
    case "set_reminder":
      //add a reminder 
      response = queueTwilio(t.Result.Parameters.DateTime, t.Result.ResolvedQuery)
    case "netflix":
      //check if its on netflix 
    default:
      response = "Sorry, I did not get what you said"
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

  http.HandleFunc("/generate_api_token", func(w http.ResponseWriter, r *http.Request) {

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
  })

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
