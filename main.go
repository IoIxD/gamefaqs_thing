package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gocolly/colly"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", gamefaqRedir)
	r.Route("/{Location:.*}", func(r chi.Router) {
		r.Get("/", gamefaqRedir)
	})

	s := &http.Server{
		Addr:           ":8085",
		Handler:        http.HandlerFunc(gamefaqRedir),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

type ArchivedResult struct {
	ArchivedSnapshots ArchivedSnapshots `json:"archived_snapshots"`
}

type ArchivedSnapshots struct {
	ArchivedSnapshot *ArchivedSnapshot `json:"closest"`
}

type ArchivedSnapshot struct {
	status    int
	avaliable bool
	url       string
	timestamp int64
}

var cachedPages map[string]CachedPage = make(map[string]CachedPage)
var knownURLs map[string]KnownURL = make(map[string]KnownURL)

type CachedPage struct {
	page   []byte
	Cached int64
}

type KnownURL struct {
	KnownAsOf int64
	Outdated  bool
}

func gamefaqRedir(w http.ResponseWriter, r *http.Request) {
	url := strings.Join(strings.Split(r.URL.Path, "/")[1:], "/")

	cached, ok := cachedPages[url]
	// if it was cacxhed
	if ok {
		// and that cache is new
		if cached.Cached >= int64(time.Now().Unix()-(int64(time.Minute)*5)) {
			w.Write(cached.page)
			return
		}

	}

	c := colly.NewCollector()

	var yes bytes.Buffer

	fmt.Println("using cached listing")

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("href"), "/sitemap") {
			link := strings.ReplaceAll(e.Attr("href"), "/sitemap", "")
			yes.Write([]byte("<a href='" + link + "'>" + link + "</a><br>"))
		} else {
			k, ok := knownURLs[e.Attr("href")]
			// if it was cacxhed
			if ok {
				// and that cache is new
				if !k.Outdated && k.KnownAsOf >= int64(time.Now().Unix()-(int64(time.Minute)*30)) {
					yes.Write([]byte("<a href='" + e.Attr("href") + "'>" + e.Attr("href") + "</a><br>"))
					return
				} else {
					k.Outdated = true
				}
			}

			res, err := http.Get("https://archive.org/wayback/available?url=" + e.Attr("href"))
			if err != nil {
				yes.Write([]byte("<em>" + err.Error() + "</em><br>"))
				return
			}

			body := bytes.NewBuffer(nil)
			io.Copy(body, res.Body)
			var r ArchivedResult
			err = json.Unmarshal(body.Bytes(), &r)
			if err != nil {
				yes.Write([]byte("<em>" + err.Error() + "</em><br>"))
				return
			}
			if r.ArchivedSnapshots.ArchivedSnapshot == nil {
				yes.Write([]byte("<a href='" + e.Attr("href") + "'>" + e.Attr("href") + "</a><br>"))
			} else {
				return
			}
			/*

				}*/

			knownURLs[url] = KnownURL{
				KnownAsOf: time.Now().Unix(),
				Outdated:  false,
			}
		}
	})
	c.Visit("https://gamefaqs.gamespot.com/sitemap/" + url)

	cachedPages[url] = CachedPage{yes.Bytes(), time.Now().Unix()}
	w.Write(yes.Bytes())
}
