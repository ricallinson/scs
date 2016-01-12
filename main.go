package main

import (
	"encoding/json"
	"fmt"
	"github.com/buger/goterm"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	_ "image/jpeg"
	"io/ioutil"
	"net/http"
	"time"
	"flag"
)

type PhotoGroup struct {
	Photos Photos `json:"photos"`
}

type Photos struct {
	Page    int8    `json:"page"`
	Pages   int8    `json:"pages"`
	PerPage int8    `json:"perpage"`
	Total   string    `json:"total"`
	Photo   []Photo `json:"photo"`
	Stat    []Photo `json:"stat"`
}

type Photo struct {
	Id       string `json:"id"`
	Owner    string `json:"owner"`
	Secret   string `json:"secret"`
	Server   string `json:"server"`
	Farm     int8   `json:"farm"`
	Title    string `json:"title"`
	IsPublic int8   `json:"ispublic"`
	IsFriend int8   `json:"isfriend"`
	IsFamily int8   `json:"isfamily"`
}

func GetFlickrImages(uid string) []string {
	images := []string{}
	apiKey := "d6d164a9fc70717b82c3d2b65847d870"
	resp, err := http.Get("https://api.flickr.com/services/rest/?method=flickr.people.getPublicPhotos&api_key=" + apiKey + "&user_id=" + uid + "&format=json&nojsoncallback=1")
	if err != nil {
		fmt.Println(err)
		return images
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return images
	}
	var photos PhotoGroup
	err = json.Unmarshal(body, &photos)
	if err != nil {
		fmt.Println(err)
		return images
	}
	for _, photo := range photos.Photos.Photo {
		images = append(images, fmt.Sprintf("https://farm%d.staticflickr.com/%s/%s_%s.jpg\n", photo.Farm, photo.Server, photo.Id, photo.Secret))
	}
	return images
}

func PrintImage(img image.Image) {
	width := goterm.Width()
	height := goterm.Height()
	img = resize.Resize(uint(width), uint(height), img, resize.NearestNeighbor)
	buf := ""
	for y := 0; y < height; y++ {
		buf += "\n"
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			grayColor := color.Gray16{Y: uint16((r + g + b) / 3)}
			pixelColor := 232 + (grayColor.Y / 255 / 16)
			buf += fmt.Sprintf("\033[38;5;#%dm█\033[m", pixelColor)
		}
	}
	goterm.Print(buf)
	goterm.Flush()
}

func PrintImageFromUrl(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	PrintImage(img)
}

func main() {
	var flickr string
	flag.StringVar(&flickr, "flickr", "", "Flickr user id") // "50566068%40N00"
	var wait int
	flag.IntVar(&wait, "wait", 5, "Number of seconds between each image") // "50566068%40N00"
	flag.Parse()
	urls := []string{}
	if len(flickr) > 0 {
		urls = append(urls, GetFlickrImages(flickr)...)
	}
	for _, url := range urls {
		PrintImageFromUrl(url)
		time.Sleep(time.Duration(wait) * time.Second)
	}
}
