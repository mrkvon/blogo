package main

import (
  "net/http"
  "fmt"
  "os"
  "io/ioutil"
  "regexp"
  "encoding/json"
)

var validArchivePath = regexp.MustCompile("^/archive/?$")
var validHomePath = regexp.MustCompile("^/?$")
var validPostPath = regexp.MustCompile("^/([a-z0-9-]+)/?$")

func homeHandler(w http.ResponseWriter, r *http.Request) {
  // read the json list of posts

  // match home and do accordingly
  if validHomePath.MatchString(r.URL.Path) {
    list, err:=loadList()
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    http.Redirect(w, r, "/"+list[0].Url, http.StatusFound)
    return
  }
  m:=validPostPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
    return
  }

  // match
  list, err:=loadList()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  found := false
  index :=-1
  for i, item:= range list {
    if item.Url == m[1] {
      found = true
      index = i
      break;
    }
  }

  if found {
    content, err := loadPost(list[index].Url)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    fmt.Fprint(w, content)
    return
  }

  http.NotFound(w,r)
  return
}

func archiveHandler(w http.ResponseWriter, r *http.Request) {
  if !validArchivePath.MatchString(r.URL.Path) {
    http.NotFound(w, r)
    return
  }
  list, err:=loadList()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  fmt.Fprintf(w, "<ul>\n")
  for _, item := range list {
    fmt.Fprintf(w, "<li><a href=\"/%s\">%s</a></li>\n", item.Url, item.Title)
  }
  fmt.Fprintf(w, "</ul>")
}

func main() {
  http.HandleFunc("/archive", archiveHandler)
  http.HandleFunc("/archive/", archiveHandler)
  http.HandleFunc("/", homeHandler)

  http.ListenAndServe(":8106", nil)
}

type PostItem struct {
  Date string `json:"date"`
  Title string `json:"title"`
  Url string `json:"url"`
}

// function to load the json list of posts
func loadList() ([]*PostItem, error) {
  posts, err := os.Open("posts.json")
  if err != nil {
    return nil, err
  }
  var list []*PostItem
  json.NewDecoder(posts).Decode(&list)
  return list, nil
}

func loadPost(url string) (string, error) {
  filename := "./posts/"+url+".md"
  rawContent, err:= ioutil.ReadFile(filename)
  if err != nil {
    return "", err
  }
  return string(rawContent), nil
}

