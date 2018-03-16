package nct

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ndphu/music-downloader/provider"
	"github.com/ndphu/music-downloader/utils"
	iohelper "github.com/ndphu/music-downloader/utils/io"
	"net/url"
	"os"
	"path"
	"sync"
)

func downloadAlbumOrTrack(c *provider.DownloadContext) error {
	doc, err := goquery.NewDocument(c.URL.String())
	if err != nil {
		return err
	}
	ajaxUrl, err := getAjaxUrl(doc)
	if err != nil {
		return err
	}
	pageTitle, err := getPageTitle(doc)
	if err != nil {
		return err
	}

	if ajaxUrl.String() == "" {
		return errors.New("Fail to get AjaxURL")
	} else {
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
	}

}

func downloadAlbum(trackList *Tracklist, c *provider.DownloadContext) error {
	fmt.Printf("Downloading album [%s]...\n", trackList.PageTitle)
	fmt.Printf("Found %d items\n", len(trackList.Tracks))
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
		fmt.Printf("[%d] '%s'\n", i, utils.TrimTitle(track.Title))

		w.Add(1)
		runningThread++
		go func(_t Track, _c *provider.DownloadContext) {
			defer w.Done()
			err := downloadTrack(&_t, _c)
			if err != nil {
				fmt.Printf("Fail to download track %s at %s or %s\n", _t.Title, _t.Location, _t.LocationHQ)
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

func downloadTrack(t *Track, c *provider.DownloadContext) error {
	title := iohelper.CleanupFileName(utils.TrimTitle(t.Title))
	filePath := path.Join(c.Output, title+".mp3")
	location := utils.TrimCDATA(t.Location)
	locationHQ := utils.TrimCDATA(t.LocationHQ)

	locationLossless := locationLossless(t)
	if locationLossless != "" {
		fmt.Println("Downloading song " + title + " (Lossless)...")
		fmt.Println("URL = " + locationLossless)
		return iohelper.DownloadFileWithRetry(filePath, locationLossless, 5)
	} else if locationHQ != "" {
		fmt.Println("Downloading song " + title + " (VIP)...")
		return iohelper.DownloadFileWithRetry(filePath, locationHQ, 5)
	} else {
		fmt.Println("Downloading song " + title + "...")
		return iohelper.DownloadFileWithRetry(filePath, location, 5)
	}

}

func locationLossless(t *Track) string {
	losslessUrl, _ := url.Parse(fmt.Sprintf("https://www.nhaccuatui.com/download/song/%s_lossless", utils.TrimCDATA(t.Key)))
	savedCookie := getAuthCookie()
	var data []byte
	var err error

	if savedCookie == nil {
		return ""
	} else {
		headers := make(map[string]string)
		headers["X-Requested-With"] = "XMLHttpRequest"
		headers["Referer"] = utils.TrimCDATA(t.Info)
		data, err = iohelper.GetWithCookie(losslessUrl, savedCookie, headers)
	}
	if err != nil {
		panic(err)
		return ""
	}
	losslessResponse := LosslessResponse{}

	err = json.Unmarshal(data, &losslessResponse)
	if err != nil {
		return ""
	}
	if losslessResponse.Data["is_charge"] == "false" {
		return ""
	}
	return losslessResponse.Data["stream_url"]
}
