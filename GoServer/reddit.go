package main 



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

func getRedditTopPost() string {
  response := RedditRes{}
  simple_req("GET", "https://www.reddit.com/r/all.json", nil, nil, &response)
  return response.Data.Children[0].Data.Title
}


//****************** END reddit **********************//