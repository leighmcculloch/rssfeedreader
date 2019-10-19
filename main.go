package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/mmcdole/gofeed"
)

const maxItems = 15

type feed struct {
	Name string
	URL  string
}

var feedDefs = []feed{
	{"SF Gate", "https://www.sfgate.com/bayarea/feed/Bay-Area-News-429.php"},
	{"ABC", "https://www.abc.net.au/news/feed/51120/rss.xml"},
	{"Supercars", "https://www.supercars.com/rss/news.rss"},
	{"NY Times", "https://www.nytimes.com/services/xml/rss/nyt/HomePage.xml"},
	{"HN", "http://news.ycombinator.com/rss"},
	{"Lobsters", "https://lobste.rs/rss"},
	{"Go News", "https://golangnews.com/index.xml"},
	{"NPR", "https://www.npr.org/rss/rss.php?id=1002"},
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
		if len(f.Items) > maxItems {
			f.Items = f.Items[:maxItems]
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
	err := indexTemplate().Execute(w, d)
	if err != nil {
		log.Printf("template error: %v", err)
	}
}

func indexTemplate() *template.Template {
	return template.Must(template.ParseFiles("index.html"))
}

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
