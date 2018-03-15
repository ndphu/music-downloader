package nct

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ndphu/music-downloader/provider"
	"net/url"
)

func downloadAllAlbumForArtist(c *provider.DownloadContext) error {
	current := c.URL

	doc, err := goquery.NewDocument(current.String())
	if err != nil {
		return err
	}
	var albums []*url.URL
	albums = append(albums, parseAlbumInPage(doc)...)

	for {
		if isLastPage(doc) {
			break
		}
		current = getNextPage(doc)
		if current == nil {
			break
		}
		doc, err = goquery.NewDocument(current.String())
		if err != nil {
			return err
		}
		albums = append(albums, parseAlbumInPage(doc)...)
	}

	fmt.Printf("Found %d album(s)\n", len(albums))
	for _, album := range albums {
		albumContext := &provider.DownloadContext{
			URL:         album,
			Indexes:     c.Indexes,
			Output:      c.Output,
			ThreadCount: c.ThreadCount,
		}
		err = downloadAlbumOrTrack(albumContext)
		if err != nil {
			return err
		}
	}

	return nil
}

func isLastPage(doc *goquery.Document) bool {
	return doc.Find(".box_pageview a").Last().HasClass("active")
}

func getNextPage(doc *goquery.Document) *url.URL {
	var next *url.URL
	foundCurrent := false
	doc.Find(".box_pageview a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if foundCurrent {
			// this is next page
			next, _ = url.Parse(s.AttrOr("href", ""))
			return false
		} else {
			foundCurrent = s.HasClass("active")
			return true
		}
	})
	return next
}

func parseAlbumInPage(doc *goquery.Document) []*url.URL {
	var albums []*url.URL
	doc.Find(".list_album .info_album h3 a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		href, err := url.Parse(s.AttrOr("href", ""))
		albums = append(albums, href)
		return err == nil
	})
	return albums
}
