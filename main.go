package main

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/mmcdole/gofeed"
)

type feed struct {
	Name     string
	URL      string
	MaxItems int
}

var feedDefs = []feed{
	{"SF Gate", "https://www.sfgate.com/bayarea/feed/Bay-Area-News-429.php", 6},
	{"ABC", "https://www.abc.net.au/news/feed/51120/rss.xml", 15},
	{"Supercars", "https://www.supercars.com/rss/news.rss", 12},
	{"Public Node", "https://publicnode.org/feed/", 6},
	{"HN", "http://news.ycombinator.com/rss", 15},
	{"Lobsters", "https://lobste.rs/rss", 15},
	{"Go News", "https://golangnews.com/index.xml", 7},
	{"NPR", "https://www.npr.org/rss/rss.php?id=1002", 10},
	{"NY Times", "https://www.nytimes.com/services/xml/rss/nyt/HomePage.xml", 15},
}

func index(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Feeds []*gofeed.Feed
	}

	d := data{}
	for _, fd := range feedDefs {
		fp := gofeed.NewParser()
		f, err := fp.ParseURL(fd.URL)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		if len(f.Items) > fd.MaxItems {
			f.Items = f.Items[:fd.MaxItems]
		}
		if fd.Name != "" {
			f.Title = fd.Name
		}
		if fd.Name == "Plant Eater" {
			for _, i := range f.Items {
				i.Title = i.Description
			}
		}
		d.Feeds = append(d.Feeds, f)
	}
	err := indexTemplate.Execute(w, d)
	if err != nil {
		log.Printf("template error: %v", err)
	}
}

//go:embed index.html
var indexHTML string

var indexTemplate = func() *template.Template {
	return template.Must(template.New("index.html").Parse(indexHTML))
}()

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
