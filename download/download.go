package download

import (
	"github.com/ndphu/music-downloader/download/zing"
	"gopkg.in/urfave/cli.v2"
	"log"
	"net/url"
	"os"
)

type Downloader interface {
	GetSupportedSites() []string
	IsSiteSupported(string) bool
	Download(inputUrl *url.URL, outputDir string) error
}

var downloaders = []Downloader{zing.NewZingDownloader()}

func HandleDownload(c *cli.Context) error {
	rawurl := c.Args().Get(0)
	log.Printf("Input URL = %s\n", rawurl)

	outputDir := c.String("output")
	log.Printf("Output directory = %s\n", outputDir)

	err := os.MkdirAll(outputDir, 0777)
	if err != nil {
		log.Panic(err)
	}

	inputUrl, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	hostname := inputUrl.Hostname()
	log.Printf("Downloading from host: %s\n", hostname)
	for _, downloader := range downloaders {
		if downloader.IsSiteSupported(hostname) {
			err = downloader.Download(inputUrl, outputDir)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
