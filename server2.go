package main

import (
	"fmt"
	"net/http"
	"encoding/json"
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

func main() {
	http.HandleFunc("/api", handler)
	http.ListenAndServe(":8080", nil)
}

