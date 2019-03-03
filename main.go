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
	Name           string
	URL            string
	UseDescription bool
}

var feedDefs = []feed{
	{"SF Gate", "https://www.sfgate.com/bayarea/feed/Bay-Area-News-429.php", false},
	{"ABC", "https://www.abc.net.au/news/feed/51120/rss.xml", false},
	{"NY Times", "https://www.nytimes.com/services/xml/rss/nyt/HomePage.xml", false},
	{"Plant Eater", "https://web.stagram.com/rss/n/planet_eater", true},
	{"HN", "http://news.ycombinator.com/rss", false},
	{"Lobsters", "https://lobste.rs/rss", false},
	{"Go News", "https://golangnews.com/index.xml", false},
	{"NPR", "https://www.npr.org/rss/rss.php?id=1002", false},
	{"Washington Post Politics", "http://feeds.washingtonpost.com/rss/politics", false},
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
