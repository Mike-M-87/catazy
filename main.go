package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Noty struct {
	Quality  string
	Title    string
	Year     string
	Image    string
	FilmType string
	WatchUrl string
	Duration string
}

func main() {
	Catazz()
}

func Catazz() ([]Noty, error) {
	resp, err := http.DefaultClient.Get("https://cataz.to/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	reBody := regexp.MustCompile(`<div class="flw-item">[\s\S]*?<div class="clearfix"></div>(\n)</div>`)
	notIds := reBody.FindAllString(string(responseBody), -1)

	titleRe := regexp.MustCompile(`<h3 class="film-name"><a\s+[^>]+>([^<]+)</a>\s+</h3>`)
	qualityRe := regexp.MustCompile(`<div class="pick film-poster-quality">(.*?)</div>`)
	imgRe := regexp.MustCompile(`<img data-src="([^"]+)"`)
	yearRe := regexp.MustCompile(`<span class="fdi-item">(\d{4})</span>`)
	typeRe := regexp.MustCompile(`<span class="float-right fdi-type">([^<]+)</span>`)
	minsRe := regexp.MustCompile(`<span class="fdi-item fdi-duration">([^<]+)</span>`)
	urlRe := regexp.MustCompile(`<a\s+href="([^"]+)"\s+class="[^"]+"\s+`)

	var noties []Noty
	fmt.Println("🟡🟡🟡🟡🟡🟡🟡")

	for _, v := range notIds {
		fmt.Print()
		fpq := qualityRe.FindStringSubmatch(v)
		tt := titleRe.FindStringSubmatch(v)
		img := imgRe.FindStringSubmatch(v)
		year := yearRe.FindStringSubmatch(v)
		minutes := minsRe.FindStringSubmatch(v)
		filmType := typeRe.FindStringSubmatch(v)
		watchurl := urlRe.FindStringSubmatch(v)

		var noty Noty

		if len(fpq) > 1 && len(tt) > 1 && len(img) > 1 && len(year) > 1 && len(filmType) > 1 && len(minutes) > 1 && len(watchurl) > 1 {
			noty.Quality = fpq[1]
			noty.Title = tt[1]
			noty.Duration = minutes[1]
			noty.Image = img[1]
			noty.FilmType = filmType[1]
			noty.WatchUrl = watchurl[1]
			noty.Year = year[1]
			noties = append(noties, noty)
		} else {
			fmt.Println(v, "❇️", fpq, "❎", tt, "❎", img, "❎", year, "❎", minutes, "❎", filmType, "❎", watchurl, "❎", "❌❌❌\n\n")
		}
	}

	fmt.Println(len(noties), len(notIds))
	return noties, nil
}

func AsJson(input interface{}) string {
	jsonB, _ := json.Marshal(input)
	return string(jsonB)
}

func AsPrettyJson(input interface{}) string {
	jsonB, _ := json.MarshalIndent(input, "", "  ")
	return string(jsonB)
}
