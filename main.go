package main

import (
	"github.com/ndphu/music-downloader/download"
	"github.com/ndphu/music-downloader/download/nct"
	"github.com/ndphu/music-downloader/download/zing"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
)

func main() {

	downloadHandler := download.NewDownloadHandler([]download.Downloader{
		zing.NewZingDownloader(),
		nct.NewDownloader(),
	})

	app := &cli.App{
		Name:    "music-downloader",
		Usage:   "download music from multiple sources",
		Version: "0.0.2",
		Commands: []*cli.Command{
			{
				Name:    "download",
				Aliases: []string{"dl"},
				Usage:   "Download music file(s) from a URL.\nInput could be a link to a playlist (album) or a single song",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o", "out"},
						Value:   ".",
						Usage:   "Specify the output location to save downloaded file(s)",
					},
					&cli.StringFlag{
						Name:    "index",
						Aliases: []string{"i"},
						Value:   "",
						Usage:   "List for song tobe downloaded (by index)",
					},
					&cli.IntFlag{
						Name:    "thread-count",
						Aliases: []string{"n"},
						Value:   1,
						Usage:   "Number of parallel download. No parallel download by default.",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						log.Panic("Input URL(s) is required")
					}
					err := downloadHandler.HandleDownload(c)
					if err != nil {
						panic(err)
					} else {
						log.Println("Done without error")
					}
					return err
				},
			},
			{
				Name:  "login",
				Usage: "Login to the music provider to download high quality files",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Value: "",
						Usage: "The the login host",
					},
					&cli.StringFlag{
						Name:    "username",
						Aliases: []string{"u"},
						Value:   "",
						Usage:   "The login username",
					},
					&cli.StringFlag{
						Name:    "password",
						Aliases: []string{"p"},
						Value:   "",
						Usage:   "The login password",
					},
				},
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}
	app.Run(os.Args)
}
