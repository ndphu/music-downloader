package zing

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ndphu/music-downloader/provider"
	iohelper "github.com/ndphu/music-downloader/utils/io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	ZING_HOST string = "mp3.zing.vn"
)

type ZingProvider struct {
}

func NewProvider() *ZingProvider {
	return &ZingProvider{}
}

func (*ZingProvider) GetName() string {
	return "zing"
}

func (*ZingProvider) GetSupportedSites() []string {
	return []string{ZING_HOST}
}

func (*ZingProvider) IsSiteSupported(site string) bool {
	return site == ZING_HOST
}

func (d *ZingProvider) Download(context *provider.DownloadContext) error {
	doc, err := goquery.NewDocument(context.URL.String())
	if err != nil {
		return err
	}

	playerElement := doc.Find("div#zplayerjs-wrapper").First()
	dataXML := playerElement.AttrOr("data-xml", "")
	if dataXML == "" {
		return errors.New("Fail to get data-xml for media source")
	}

	return d.download(doc, dataXML, context)
}

func (d *ZingProvider) download(doc *goquery.Document, dataXML string, context *provider.DownloadContext) error {
	log.Printf("Data XML = %s\n", dataXML)
	dataUrl, err := url.Parse(fmt.Sprintf("https://%s/xhr%s", ZING_HOST, dataXML))

	if err != nil {
		return err
	}

	itemType := dataUrl.Query().Get("type")

	if itemType == "audio" {
		return downloadSong(dataUrl, context)
	} else {
		return downloadAlbum(doc, dataUrl, context)
	}
}

func downloadSong(dataUrl *url.URL, context *provider.DownloadContext) error {
	resp, err := http.Get(dataUrl.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("Reading body...")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("Body size = %d\n", len(body))

	songResponse := SongResponse{}
	err = json.Unmarshal(body, &songResponse)
	if err != nil {
		return err
	}

	if songResponse.Err != 0 {
		return errors.New(songResponse.Msg)
	}

	return iohelper.DownloadFileWithRetry(songResponse.Data.Name+".mp3", "https:"+songResponse.Data.Source.Normal, 5)
}

func downloadAlbum(doc *goquery.Document, dataUrl *url.URL, context *provider.DownloadContext) error {
	title := strings.Trim(doc.Find("div.info-content h1").First().Text(), " ")
	log.Println("Album tile is \"" + title + "\"")
	outputDir := context.Output + "/" + iohelper.CleanupFileName(title)
	err := os.MkdirAll(outputDir, 0777)
	if err != nil {
		log.Panic(err)
	}

	resp, err := http.Get(dataUrl.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("Reading body...")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("Body size = %d\n", len(body))

	albumData := AlbumResponse{}
	err = json.Unmarshal(body, &albumData)
	if err != nil {
		return err
	}

	if albumData.Err != 0 {
		return errors.New(albumData.Msg)
	}

	numOfItem := len(albumData.Data.Items)

	log.Printf("Found %d item(s)\n", numOfItem)
	filterByIndex := len(context.Indexes) > 0

	if filterByIndex {
		log.Printf("Only download song(s) with index = %v\n", context.Indexes)
	}

	w := sync.WaitGroup{}
	runningThreadCount := 0

	for _, item := range albumData.Data.Items {
		log.Printf("%v\n", item.Order)
		if filterByIndex {
			itemOrder := -1
			switch item.Order.(type) {
			case string:
				itemOrder, err = strconv.Atoi(item.Order.(string))
				if err != nil {
					return err
				}
			case float64:
				float, _ := item.Order.(float64)
				itemOrder = int(float)
			}

			include, err := contains(context.Indexes, itemOrder)
			if err != nil {
				return err
			}
			if !include {
				continue
			}
		}

		item.Name = iohelper.CleanupFileName(item.Name)

		log.Printf("Downloading item \"%s\"\n", item.Name)
		var downloadUrl string
		if item.IsVip {
			downloadUrl = item.Source.High
		} else {
			downloadUrl = item.Source.Normal
		}

		if strings.Index(downloadUrl, "//") == 0 {
			downloadUrl = "https:" + downloadUrl
		}
		w.Add(1)
		runningThreadCount++
		go func(_url string, _name string) {

			defer w.Done()
			err = iohelper.DownloadFileWithRetry(fmt.Sprintf("%s/%s.mp3", outputDir, _name), _url, 5)
			if err != nil {
				log.Panic(err)
			}
		}(downloadUrl, item.Name)
		if runningThreadCount == context.ThreadCount {
			w.Wait()
			runningThreadCount = 0
		}
	}

	w.Wait()

	return nil
}

func contains(arr []int, val int) (bool, error) {
	for _, cur := range arr {
		if val == cur {
			return true, nil
		}
	}
	return false, nil
}
