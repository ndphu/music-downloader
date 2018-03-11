package download

import (
	"gopkg.in/urfave/cli.v2"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Downloader interface {
	GetSupportedSites() []string
	IsSiteSupported(string) bool
	Download(*DownloadContext) error
}

type DownloadContext struct {
	URL         *url.URL
	Output      string
	Indexes     []int
	ThreadCount int
}

type DownloadHandler struct {
	downloaders []Downloader
}

func NewDownloadHandler(_downloaders []Downloader) *DownloadHandler {
	handler := &DownloadHandler{
		downloaders: _downloaders,
	}
	return handler
}

func (handler *DownloadHandler) HandleDownload(c *cli.Context) error {

	outputDir := c.String("output")
	log.Printf("Output directory = %s\n", outputDir)

	var indexInt []int
	option := c.String("index")

	if option == "" {
		indexInt = make([]int, 0)
	} else {
		indexes := strings.Split(option, ",")
		indexInt = make([]int, len(indexes))
		for i, idxStr := range indexes {
			intVal, err := strconv.Atoi(idxStr)
			if err != nil {
				return err
			}
			indexInt[i] = intVal
		}
	}

	err := os.MkdirAll(outputDir, 0777)
	if err != nil {
		log.Panic(err)
	}

	for i := 0; i < c.Args().Len(); i++ {
		rawurl := c.Args().Get(i)
		log.Printf("Input URL = %s\n", rawurl)

		inputUrl, err := url.Parse(rawurl)
		if err != nil {
			return err
		}

		context := &DownloadContext{
			URL:         inputUrl,
			Output:      outputDir,
			Indexes:     indexInt,
			ThreadCount: c.Int("thread-count"),
		}

		hostname := inputUrl.Hostname()
		log.Printf("Downloading from host: %s\n", hostname)
		for _, h := range handler.downloaders {
			if h.IsSiteSupported(hostname) {
				err = h.Download(context)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
