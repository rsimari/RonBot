package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
  "io/ioutil"
	//"net/http/httputil"
)

type WebhookResult struct {
	Action string
  Fulfillment WebhookFulfillment
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

func getRedditTopPost() string {
  response := RedditRes{}
  simple_req("GET", "https://www.reddit.com/r/all.json", nil, nil, &response)
  return response.Data.Children[0].Data.Title
}


//handles making a simply request 
func simple_req(method string, url string, headers map[string]string, body map[string]string, response interface{}) error {

  var client http.Client
  r, err := http.NewRequest(method, url, nil)

  for k, v := range headers {
    r.Header.Add(k, v)
  }

  if body != nil {
    jsonString, err := json.Marshal(body)
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
  var response string
  // reads action from api.ai bot and sends back correct api response
  switch t.Result.Action {
    case "hacker_news_top_post":
      response = "The top post of hacker news is: " + getHackerNewsTopPost()
    case "hacker_news_top_job":
      response = "The top job post on hacker news is currently: " + getHackerNewsTopJob()
    case "reddit_top_post":
      response = "The top post on reddit right now is: " + getRedditTopPost()
    default:
      response = "Sorry, I didnt get what you said"
  }

	fmt.Fprintf(rw, "{ \"speech\": \"%s\" }", response)
}

func main() {
  //textHandler := http.HandlerFunc(textPost)
	//http.Handle("/api/text", authorization(textHandler))
	http.HandleFunc("/api", webhook_handler)

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
