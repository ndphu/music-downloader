package provider

import (
	"errors"
	"fmt"
	"gopkg.in/urfave/cli.v2"

	"net/url"
	"os"
	"strconv"
	"strings"
)

type DownloadContext struct {
	URL         *url.URL
	Output      string
	Indexes     []int
	ThreadCount int
}

func (providerService *ProviderService) HandleDownload(c *cli.Context) error {

	outputDir := c.String("output")
	fmt.Printf("Output directory = %s\n", outputDir)

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
		panic(err)
	}

	for i := 0; i < c.Args().Len(); i++ {
		rawurl := c.Args().Get(i)

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

		hadProvider := false
		hostname := inputUrl.Hostname()

		for _, p := range providerService.providers {
			if p.IsSiteSupported(hostname) {
				hadProvider = true
				err = p.Download(context)
				if err != nil {
					return err
				}
			}
		}
		if !hadProvider {
			return errors.New(fmt.Sprintf("Current music provider [%s] is not supported", hostname))
		}
	}

	return nil
}
