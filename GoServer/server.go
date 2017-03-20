package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
  "io/ioutil"
  "bytes"
	"strings"

	//"net/http/httputil"
)


//************** IRON.IO STUFF *******************//

type Task struct {
        CodeName string `json:"code_name"`
        Payload string `json:"payload"`
        StartAt string `json:"start_at"`
}

type ReqData struct {
        Tasks []*Task `json:"tasks"`
}

func queueTwilio(execTime string, msg string) string {
        const token = "IpsW3ZNlFsLl42T2vSqw"
        const project = "58cf2ffd0e2c7300061ad812"

        // Insert our project ID and token into the API endpoint
        target := fmt.Sprintf("http://worker-us-east.iron.io/2/projects/%s/schedules?oauth=%s", project, token)

        // Build the payload
        // The payload is a string to pass information into your worker as part of a task
        // It generally is a JSON-serialized string (which is what we're doing here) that can be deserialized in the worker
        payload := map[string]interface{} {
                "Message" : "Test",
        }

        payload_bytes, err := json.Marshal(payload)
        if err != nil {
                panic(err.Error())
        }
        payload_str := string(payload_bytes)

        // Build the task
        task := &Task {
                CodeName: "irontest",
                Payload: payload_str,
                StartAt: execTime,
        }

        // Build a request containing the task
        json_data := &ReqData {
                Tasks: []*Task { task },
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

  return "It is " + description + " in " + city + " right now. The temperature in is " + s + " degrees right now, with a high of " + sHigh + " and a low of " + sLow + "."
}




//handles making a simply request 
func simple_req(method string, url string, headers map[string]string, body map[string]string, response interface{}) error {

  var client http.Client

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
    case "current_time": 
      //get the current time 
    case "reminder":
      //add a reminder 
      response = queueTwilio(t.Result.Parameters.DateTime, t.Result.ResolvedQuery)
    case "netflix": 
      //check if its on netflix 
    default:
      response = "Sorry, I didnt get what you said"
  }

	fmt.Fprintf(rw, "{ \"speech\": \"%s\" }", response)
}



func main() {
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
