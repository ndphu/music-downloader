package zing

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	ZING_HOST string = "mp3.zing.vn"
)

type ZingDownloader struct {
}

func NewZingDownloader() *ZingDownloader {
	return &ZingDownloader{}
}

func (*ZingDownloader) GetSupportedSites() []string {
	return []string{ZING_HOST}
}

func (*ZingDownloader) IsSiteSupported(site string) bool {
	return site == ZING_HOST
}

func (d *ZingDownloader) Download(input *url.URL, outputDir string) error {
	doc, err := goquery.NewDocument(input.String())
	if err != nil {
		return err
	}

	playerElement := doc.Find("div#zplayerjs-wrapper").First()
	dataXML := playerElement.AttrOr("data-xml", "")
	if dataXML == "" {
		return errors.New("Fail to get data-xml for media source")
	}

	return d.download(doc, dataXML, outputDir)
}

func (d *ZingDownloader) download(doc *goquery.Document, dataXML string, outputDir string) error {
	log.Printf("Data XML = %s\n", dataXML)
	dataUrl, err := url.Parse(fmt.Sprintf("https://%s/xhr%s", ZING_HOST, dataXML))

	if err != nil {
		return err
	}

	itemType := dataUrl.Query().Get("type")

	if itemType == "audio" {
		return downloadSong(dataUrl)
	} else {
		return downloadAlbum(doc, dataUrl, outputDir)
	}
}

func downloadSong(dataUrl *url.URL) error {
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

	return downloadFileWithRetry(songResponse.Data.Name+".mp3", "https:"+songResponse.Data.Source.Normal, 5)
}

func downloadAlbum(doc *goquery.Document, dataUrl *url.URL, outputDir string) error {
	title := doc.Find("div.info-content h1").First().Text()
	outputDir = outputDir + "/" + title
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

	for _, item := range albumData.Data.Items {
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

		err = downloadFileWithRetry(fmt.Sprintf("%s/%s.mp3", outputDir, item.Name), downloadUrl, 5)
		if err != nil {
			log.Panic(err)
			return err
		}
	}

	return nil
}

func downloadFileWithRetry(filepath string, fileUrl string, retry int) (err error) {
	try := 0

	for {
		try++
		err = downloadFile(filepath, fileUrl)
		if err == nil {
			return err
		}
		if try == retry {
			return err
		} else {
			log.Printf("%v\n", err)
			log.Printf("Retrying... %d\n", try)
		}
	}
	return err
}

func downloadFile(filepath string, fileUrl string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
