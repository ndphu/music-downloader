package nct

import (
	"encoding/xml"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/ndphu/music-downloader/download"
	"github.com/ndphu/music-downloader/utils"
	iohelper "github.com/ndphu/music-downloader/utils/io"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
)

var (
	NCT_HOST = "www.nhaccuatui.com"
)

type NCTDownloader struct {
}

func NewDownloader() *NCTDownloader {
	return &NCTDownloader{}
}

func (*NCTDownloader) GetSupportedSites() []string {
	return []string{NCT_HOST}
}

func (*NCTDownloader) IsSiteSupported(site string) bool {
	return site == NCT_HOST
}

func (*NCTDownloader) Download(c *download.DownloadContext) error {
	ajaxUrl, pageTitle, err := crawWebPage(c.URL)
	if err != nil {
		return err
	}
	log.Println("AjaxURL = " + ajaxUrl.String())

	trackList, err := getTracklist(ajaxUrl)
	if err != nil {
		return err
	}

	trackList.PageTitle = utils.TrimTitle(pageTitle)

	if trackList.Type == "song" {
		return downloadTrack(&trackList.Tracks[0], c)
	} else if trackList.Type == "playlist" {
		return downloadAlbum(trackList, c)
	} else {
		return errors.New("Unsupported URL with type = " + trackList.Type)
	}

	return nil
}

func downloadAlbum(trackList *Tracklist, c *download.DownloadContext) error {
	log.Printf("Found %d items\n", len(trackList.Tracks))
	c.Output = path.Join(c.Output, trackList.PageTitle)
	err := os.MkdirAll(c.Output, 0777)
	if err != nil {
		panic(err)
	}

	w := sync.WaitGroup{}
	runningThread := 0
	for i, track := range trackList.Tracks {
		if len(c.Indexes) > 0 && !utils.ArrayContains(c.Indexes, i+1) {
			continue
		}
		log.Printf("[%d] '%s'\n", i, utils.TrimTitle(track.Title))

		w.Add(1)
		runningThread++
		go func(_t Track, _c *download.DownloadContext) {
			defer w.Done()
			err := downloadTrack(&_t, _c)
			if err != nil {
				panic(err)
			}
		}(track, c)
		if runningThread == c.ThreadCount {
			w.Wait()
			runningThread = 0
		}
	}

	w.Wait()
	return nil
}

func downloadTrack(t *Track, c *download.DownloadContext) error {
	title := utils.TrimTitle(t.Title)
	filePath := iohelper.CleanupFileName(path.Join(c.Output, title+".mp3"))
	log.Println("Downloading song " + title + "...")
	return iohelper.DownloadFileWithRetry(filePath, utils.TrimTitle(t.Location), 5)
}

func crawWebPage(input *url.URL) (ajaxUrl *url.URL, title string, err error) {
	doc, err := goquery.NewDocument(input.String())
	if err != nil {
		return nil, "", err
	}

	rawAjaxUrl := ""

	doc.Find("div.playing_absolute script").EachWithBreak(func(i int, s *goquery.Selection) bool {
		script := s.Text()
		line := strings.Split(script, "\n")
		for _, line := range line {
			if strings.Index(line, "player.peConfig.xmlURL") > 0 {
				rawAjaxUrl = strings.Split(line, "\"")[1]
			}
		}
		return rawAjaxUrl == ""
	})

	title = doc.Find("title").First().Text()

	ajaxUrl, err = url.Parse(rawAjaxUrl)

	return ajaxUrl, title, err
}

func getTracklist(ajaxUrl *url.URL) (*Tracklist, error) {
	data, err := iohelper.ReadFromUrl(ajaxUrl)
	if err != nil {
		return nil, err
	}

	log.Printf("Data size = %d\n", len(data))

	resp := Tracklist{}

	err = xml.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
