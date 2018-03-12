package main

import (
	"fmt"
	"github.com/ndphu/music-downloader/provider"
	"github.com/ndphu/music-downloader/provider/nct"
	"github.com/ndphu/music-downloader/provider/zing"
	"gopkg.in/urfave/cli.v2"
	"os"
)

func main() {

	providerService := provider.NewProviderService([]provider.Provider{
		zing.NewProvider(),
		nct.NewProvider(),
	})

	app := &cli.App{
		Name:    "music-downloader",
		Usage:   "download music from multiple sources",
		Version: "0.0.4",
		Commands: []*cli.Command{
			{
				Name:      "download",
				Aliases:   []string{"dl"},
				Usage:     "Download music file(s) from a URL.\nInput could be a link to a playlist (album) or a single song",
				Flags:     provider.DownloadFlags,
				ArgsUsage: "[link_1] [link_2]...[link_n]",
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						fmt.Println("Need to provide input URL(s)")
					}
					err := providerService.HandleDownload(c)
					if err != nil {
						panic(err)
					} else {
						fmt.Println("Done")
					}
					return err
				},
			},
			{
				Name:  "provider",
				Usage: "Action for providers",
				Subcommands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"ls"},
						Usage:   "Shows all supported providers",
						Action:  providerService.ListProviderHandler,
					},
					{
						Name:   "login",
						Usage:  "Login to a specify provider",
						Flags:  provider.LoginFlags,
						Action: providerService.LoginActionHandler,
					},
				},
			},
		},
	}
	app.Run(os.Args)
}
