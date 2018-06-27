package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

var chosenURLs = map[int]string{}

var formTemplate = template.Must(template.New("").Parse(`
<h1><a href="{{.search}}">{{.query}}</a></h1>
<ul>
{{range .links}}
<li>
<form method="POST">
<h4><a href="{{.URL}}">{{.Title}}</a></h4>
<input type="image" src="{{.Src}}" />
<input type="hidden" name="url" value="{{.URL}}" />
</form>
</li>
{{end}}
<form method="POST">
<button type="submit">or NONE OF THE ABOVE"</button>
<input type="hidden" name="url" value="" />
</form>
</ul>
`))

func main() {
	searchQueries, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(searchQueries), "\n")
	if len(lines) > 0 {
		startServer(lines)
	}
}

func startServer(lines []string) {
	r := mux.NewRouter()
	r.HandleFunc("/", redirect)
	r.HandleFunc("/{line:[0-9]+}", lineHandler(lines))
	http.Handle("/", r)
	http.ListenAndServe(":1935", r)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/0", 302)
}

func lineHandler(lines []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		lineNum := vars["line"]
		lineAsInt, _ := strconv.Atoi(lineNum)
		_, alreadyMatched := chosenURLs[lineAsInt]
		query := lines[lineAsInt]
		switch {
		case r.Method == "GET":
			search := "https://www.youtube.com/results?search_query=" + url.QueryEscape(query)
			rr, err := http.Get(search)
			if err != nil {
				return
			}
			links := getLinks(rr.Body)
			w.Header().Set("Content-Type", "text/html")
			data := map[string]interface{}{
				"links":  links,
				"query":  query,
				"search": search,
			}
			formTemplate.Execute(w, data)
		case r.Method == "POST":
			url := r.PostFormValue("url")
			if !alreadyMatched {
				if url != "" {
					fmt.Println(url)
				}
				chosenURLs[lineAsInt] = url
			}
			nextLine := lineAsInt + 1
			if len(lines) > nextLine {
				http.Redirect(w, r, "/"+strconv.Itoa(nextLine), 302)
			}
		}
	}
}

type link struct {
	URL   string
	Src   string
	Title string
}

func getLinks(r io.Reader) []link {
	var links []link
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".item-section > li").Each(func(i int, s *goquery.Selection) {
		imgEl := s.Find("[data-ytimg=\"1\"]")
		src := imgEl.AttrOr("src", "")
		if src == "" || !strings.Contains(src, "sqp") {
			src = imgEl.AttrOr("data-thumb", "")
		}
		titleEl := s.Find(".yt-lockup-title")
		url := "https://youtube.com" + titleEl.Find("a").AttrOr("href", "")
		title := titleEl.Text()
		links = append(links, link{
			URL:   url,
			Src:   src,
			Title: title,
		})
	})
	return links
}
