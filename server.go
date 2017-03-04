package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	//"net/http/httputil"
)

type WebhookResult struct {
	Action string
}

type WebhookReq struct {
	Result WebhookResult
}

type WebhookFulfillment struct {
	Speech string
	DisplayText string
	Source string
}

type WebhookRes struct {
	Fulfillment WebhookFulfillment
}

type NewTokenStruct struct {
	 FirstName string
	 LastName string
	 Email string
	 Phone string
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

func handler(rw http.ResponseWriter, request* http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	//req, err := httputil.DumpRequest(request, true)
	//fmt.Println(string(req))
	decoder := json.NewDecoder(request.Body)
	var t WebhookReq
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(t.Result.Action)
	//res := WebhookFulfillment{ Speech: t.Result.Action, DisplayText: t.Result.Action, Source: "my dick" }
	//res := WebhookRes{ WebhookFulfillment{ Speech: t.Result.Action, DisplayText: t.Result.Action } }
	//json.NewEncoder(rw).Encode(res)
	fmt.Fprintf(rw, "{ \"speech\": \"%s\" }", t.Result.Action)
}

func signToken(tokenStruct NewTokenStruct, key interface{}) (string, error) {

	// create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "firstName": tokenStruct.FirstName,
    "lastName": tokenStruct.LastName,
	})

	if out, err := token.SignedString(key); err == nil {
		fmt.Println(out)
		return out, nil
	} else {
		return "", fmt.Errorf("Error signing token: %v", err)
	}

}


func validateToken(tokenString string) error {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	    // Don't forget to validate the alg is what you expect:
	    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	    }

	    // hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	    return []byte("EB32ODSKJN234KJNDSKJSODF89N"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	    fmt.Println(claims["foo"], claims["nbf"])
	} else {
	    fmt.Println(err)
	}

	return nil
}

func main() {

	textHandler := http.HandlerFunc(textPost)
	http.Handle("/api/text", authorization(textHandler))
	http.HandleFunc("/api", handler)

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
