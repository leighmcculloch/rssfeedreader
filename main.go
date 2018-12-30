package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

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
	{"NY Times", "https://www.nytimes.com/services/xml/rss/nyt/HomePage.xml"},
	{"HN", "http://news.ycombinator.com/rss"},
	{"Lobsters", "https://lobste.rs/rss"},
	{"NPR", "https://www.npr.org/rss/rss.php?id=1002"},
	{"Washington Post Politics", "http://feeds.washingtonpost.com/rss/politics"},
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
