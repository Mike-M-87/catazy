package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Noty struct {
	Quality           string
	Title             string
	YearOrSeason      string
	Image             string
	FilmType          string
	WatchUrl          string
	DurationOrEpisode string
}

func main() {
	films, err := Catazz()
	if err != nil {
		log.Fatal(err)
	}
	println(AsPrettyJson(films))
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
	tvSeasonRe := regexp.MustCompile(`<span class="fdi-item">(.*?)</span>`)

	typeRe := regexp.MustCompile(`<span class="float-right fdi-type">([^<]+)</span>`)
	minsRe := regexp.MustCompile(`<span class="fdi-item fdi-duration">(\d+)m</span>`)
	tvEpisodeRe := regexp.MustCompile(`<span class="fdi-item">([\d.]+)</span>`)

	urlRe := regexp.MustCompile(`<a\s+href="([^"]+)"\s+class="[^"]+"\s+`)

	noties := make([]Noty, len(notIds))

	for i, v := range notIds {
		filmType := typeRe.FindStringSubmatch(v)
		var duration []string
		var yearOrSeason []string
		var minutes string
		var year string

		if len(filmType) > 1 && filmType[1] == "TV" {
			duration = tvEpisodeRe.FindStringSubmatch(v)
			if len(duration) > 1 {
				minutes = "EP " + duration[1]
			}
			yearOrSeason = tvSeasonRe.FindStringSubmatch(v)
			if len(yearOrSeason) > 1 {
				year =  yearOrSeason[1]
			}
		} else {
			duration = minsRe.FindStringSubmatch(v)
			if len(duration) > 1 {
				durationStr := duration[1]
				durationMinutes, _ := strconv.Atoi(durationStr)
				hours := durationMinutes / 60
				mins := durationMinutes % 60
				minutes = fmt.Sprintf("%dh %dm", hours, mins)
			}
			yearOrSeason = yearRe.FindStringSubmatch(v)
			if len(yearOrSeason) > 1 {
				year = yearOrSeason[1]
			}
		}

		fpq := qualityRe.FindStringSubmatch(v)
		tt := titleRe.FindStringSubmatch(v)
		img := imgRe.FindStringSubmatch(v)
		watchurl := urlRe.FindStringSubmatch(v)

		var noty Noty

		if len(fpq) > 1 {
			noty.Quality = fpq[1]
		}
		if len(tt) > 1 {
			noty.Title = tt[1]
		}
		if len(img) > 1 {
			noty.Image = img[1]
		}
		if len(filmType) > 1 {
			noty.FilmType = filmType[1]
		}
		noty.DurationOrEpisode = minutes
		noty.YearOrSeason = year

		if len(watchurl) > 1 {
			noty.WatchUrl = "https://cataz.to/" + watchurl[1]
		}
		noties[i] = noty
	}
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
