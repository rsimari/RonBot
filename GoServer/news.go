
package main 


import (
  "fmt"
)


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