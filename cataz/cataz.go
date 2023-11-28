package cataz

import (
	"cataz/models"
	"cataz/utils"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func FetchMovies(c *fiber.Ctx) error {
	var films []models.Film
	param := false
	if c.Query("watched") == "true" {
		param = true
	}
	err := utils.DB.Where("watched = ?", param).Find(&films).Error
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(films)
}

func Watch(c *fiber.Ctx) error {
	var film *models.Film
	err := utils.DB.Where("title = ?", c.Query("title")).First(&film).Error
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	film.Watched = !film.Watched
	err = utils.DB.Save(film).Error
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusAccepted)
}

func StoreMovies() {
	movies, err := GetMovies()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range movies {
		err := utils.DB.Create(&v).Error
		if err == nil {
			fmt.Println(utils.AsPrettyJson(v))
			// send notification
		} else {
			continue
		}
		time.Sleep(time.Millisecond * 250)
	}
}

func GetMovies() ([]models.Film, error) {
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
	minsRe := regexp.MustCompile(`<span class="fdi-item fdi-duration">(\d+)m</span>`)
	urlRe := regexp.MustCompile(`<a\s+href="([^"]+)"\s+class="[^"]+"\s+`)

	noties := make([]models.Film, len(notIds))

	for i, v := range notIds {
		filmType := typeRe.FindStringSubmatch(v)
		if len(filmType) > 1 && filmType[1] == "TV" {
			continue
		}
		year := yearRe.FindStringSubmatch(v)
		fpq := qualityRe.FindStringSubmatch(v)
		tt := titleRe.FindStringSubmatch(v)
		img := imgRe.FindStringSubmatch(v)
		watchurl := urlRe.FindStringSubmatch(v)
		minutes := minsRe.FindStringSubmatch(v)

		var noty models.Film

		if len(minutes) > 1 {
			durationMinutes, _ := strconv.Atoi(minutes[1])
			hours := durationMinutes / 60
			mins := durationMinutes % 60
			noty.Duration = fmt.Sprintf("%dh %dm", hours, mins)

		}
		if len(fpq) > 1 {
			noty.Quality = fpq[1]
		}
		if len(tt) > 1 {
			noty.Title = tt[1]
		}
		if len(img) > 1 {
			noty.Image = img[1]
		}
		if len(year) > 1 {
			noty.Year = year[1]
		}
		if len(watchurl) > 1 {
			noty.WatchUrl = "https://cataz.to/" + watchurl[1]
		}
		noties[i] = noty
	}
	return noties, nil
}
