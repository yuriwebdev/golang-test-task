package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"encoding/json"

	"fmt"

	parser "github.com/PuerkitoBio/goquery"
)

var jsonResponses chan []byte

type Metadata struct {
	Title string `json:"title"`
	Price   string    `json:"price"`
	Image   string    `json:"image"`
}

type ResponseData struct {
	Url      string         `json:"url"`
	Meta     Metadata       `json:"meta"`
}

var wg sync.WaitGroup

func main() {

	http.HandleFunc("/", handleResolver)

	err := http.ListenAndServe(*BIND_ADDR, nil)

	if err != nil {
		appError(err)
	}

}

func handleResolver(rw http.ResponseWriter, rq *http.Request) {

	fmt.Println("Server started")

	if rq.Method != "POST" {
		appError(errors.New("Method not allowed"))
	}

	body, err := ioutil.ReadAll(rq.Body)

	var urls []string

	if err != nil {
		appError(err)
	}

	json.Unmarshal(body, &urls)

	jsonResponses = make(chan []byte)

	wg.Add(len(urls))

	for _, url := range urls {
		go getURL(url)

	}

	var (
		data       []ResponseData
		contentLen int
	)

	go func() {
		var dt ResponseData
		for response := range jsonResponses {
			contentLen += len(response)
			json.Unmarshal(response, &dt)
			data = append(data, dt)
			wg.Done()
		}

	}()
	wg.Wait()
	close(jsonResponses)

	resp, err := json.Marshal(data)
	if err != nil {
		appError(err)
	}
	rw.Write(resp)
	return
}



func getURL(url string) []byte {

	//defer wg.Done()

	var data ResponseData

	var jsonResult []byte

	//elements := ElementsData{}

	data.Url = url
	doc, err := parser.NewDocument(url)


	if err == nil {
		doc.Find("head").Each(func(i int, s *parser.Selection) {
			pageTitle := s.Find("title").Text()
			data.Meta.Title = pageTitle
		})

		var imageExist bool


		doc.Find("#imgTagWrapperId").Each(func(i int, s *parser.Selection) {


			image,_ := s.Find("img").Attr("data-old-hires")


			if image != "" {
				imageExist = true
				data.Meta.Image = image
			}

		})
		if !imageExist {
				data.Meta.Image,_ = doc.Find("#img-canvas").Find("img").Attr("href")
			}

		price:= doc.Find("#priceblock_ourprice_row").Find("#priceblock_ourprice").Text()

		if price == "" {
			price = doc.Find("#priceblock_saleprice_row").Find("#priceblock_saleprice").Text()
		}

		data.Meta.Price = price

	}

	jsonResult, err = json.Marshal(data)

	if err != nil {
		appError(err)
	}
	jsonResponses <- jsonResult

	return jsonResult
}

func appError(err error) {
	log.Printf("%#v", err)
	os.Exit(0)
}
