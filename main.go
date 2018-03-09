package main

import (
	"github.com/ndphu/music-downloader/download"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:    "music-downloader",
		Usage:   "download music from multiple source",
		Version: "0.0.1",
		Commands: []*cli.Command{
			{
				Name:    "download",
				Aliases: []string{"dl"},
				Usage:   "download an album from a URL",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o", "out"},
						Value:   ".",
						Usage:   "specify the output location to save downloaded file(s)",
					},
				},
				Action: func(c *cli.Context) error {
					err := download.HandleDownload(c)
					if err != nil {
						panic(err)
					} else {
						log.Println("Done without error")
					}
					return err
				},
			},
		},
	}
	app.Run(os.Args)
}
