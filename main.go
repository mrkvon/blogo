package main

import (
  "net/http"
  //"fmt"
  "os"
  "io/ioutil"
  "regexp"
  "encoding/json"
  "html/template"
  "github.com/russross/blackfriday"
  //"bytes"
)

var validArchivePath = regexp.MustCompile("^/archive/?$")
var validHomePath = regexp.MustCompile("^/?$")
var validPostPath = regexp.MustCompile("^/([a-z0-9-]+)/?$")


// var templates = template.Must(template.ParseFiles("templates/menu.html"))


func homeHandler(w http.ResponseWriter, r *http.Request) {
  // read the json list of posts

  // match home and do accordingly
  if validHomePath.MatchString(r.URL.Path) {
    list, err:=loadList()
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    http.Redirect(w, r, path + "/"+list[0].Url, http.StatusFound)
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
    ////
    type Post struct {
      Body template.HTML
      Title string
      Path string
      Prev string
      Next string
    }

    var prev string
    var next string

    if index > 0 {
      next=list[index-1].Url
    }

    if index < len(list)-1 {
      prev=list[index+1].Url
    }

    tmpl := template.Must(template.New("post.html").ParseFiles("templates/post.html", "templates/menu.html"))

    p:=Post{
      Body: template.HTML(content),
      Title: list[index].Title,
      Path: path,
      Prev: prev,
      Next: next,
    }
    err2 := tmpl.ExecuteTemplate(w, "post.html", p)

    if err2 != nil {
      http.Error(w, err2.Error(), http.StatusInternalServerError)
    }
    ////
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
  ////
  type Page struct {
    PostList []*PostItem
    Path string
  }

  tmpl := template.Must(template.New("archive.html").ParseFiles("templates/archive.html", "templates/menu.html"))

  p:=Page{PostList: list, Path: path}
  err2 := tmpl.ExecuteTemplate(w, "archive.html", p)

  if err2 != nil {
    http.Error(w, err2.Error(), http.StatusInternalServerError)
  }
  ////
}

func main() {
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
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
  defer func() {
    if err := posts.Close(); err != nil {
      panic(err)
    }
  }()
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
  content:=blackfriday.MarkdownCommon(rawContent)
  return string(content), nil
}

/*
type Page struct {
  Content template.HTML //Content string
}

func main() {
  s := "<p>Hello!</p>"

  t, err := template.New("foo").Parse(`before {{.Content}} after`)

  var buf bytes.Buffer

  err = t.ExecuteTemplate(&buf, "foo", Page{Content: template.HTML(s)})
    //err = t.ExecuteTemplate(&buf, "foo", Page{Content: s})

  if err != nil {
    fmt.Println("template error:", err)
  }

  fmt.Println(string(buf.Bytes()))
}
func markDowner(args ...interface{}) template.HTML {
  s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
  return template.HTML(s)
}

tmpl := template.Must(template.New("page.html").Funcs(template.FuncMap{"markDown": markDowner}).ParseFiles("page.html"))

err := tmpl.ExecuteTemplate(w, "page.html", p)

if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
}

{{.Body | markDown}}
*/
