package nct

import (
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ndphu/music-downloader/provider"
	iohelper "github.com/ndphu/music-downloader/utils/io"
	"net/url"
	"strings"
)

var (
	NCT_HOST = "www.nhaccuatui.com"
)

type NCTProvider struct {
}

func NewProvider() *NCTProvider {
	return &NCTProvider{}
}

func (*NCTProvider) GetName() string {
	return "nct"
}

func (*NCTProvider) GetSupportedSites() []string {
	return []string{NCT_HOST}
}

func (*NCTProvider) IsSiteSupported(site string) bool {
	return site == NCT_HOST
}

func (*NCTProvider) Download(c *provider.DownloadContext) error {
	if strings.Index(c.URL.Path, "/nghe-si-") == 0 {
		return downloadAllAlbumForArtist(c)
	} else {
		return downloadAlbumOrTrack(c)
	}
	return nil
}

func getPageTitle(doc *goquery.Document) (string, error) {
	return doc.Find("title").First().Text(), nil
}

func getAjaxUrl(doc *goquery.Document) (ajaxUrl *url.URL, err error) {
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
	return url.Parse(rawAjaxUrl)
}

func getTracklist(ajaxUrl *url.URL) (*Tracklist, error) {
	savedCookie := getAuthCookie()
	var data []byte
	var err error

	if savedCookie == nil {
		data, err = iohelper.ReadFromUrl(ajaxUrl)
	} else {
		fmt.Println("Already login. Using saved cookie...")
		data, err = iohelper.GetWithCookie(ajaxUrl, savedCookie)
	}
	if err != nil {
		return nil, err
	}

	fmt.Printf("Data size = %d\n", len(data))

	resp := Tracklist{}

	err = xml.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
