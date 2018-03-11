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
		Version: "0.0.3",
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
						Usage:   "Number of parallel provider. No parallel download by default.",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						fmt.Println("Need to provide input URL")
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
						Action: func(c *cli.Context) error {
							fmt.Println("Supported providers: ")
							for _, p := range providerService.GetProviders() {
								fmt.Println("\t" + p.GetName())
							}
							return nil
						},
					},
					{
						Name:  "login",
						Usage: "Login to a specify provider",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Aliases: []string{"n"},
								Value:   "",
								Usage:   "The name of the provider. Use 'provider ls' for all supported providers",
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
							err := providerService.Login(c)
							if err != nil {
								panic(err)
							} else {
								fmt.Println("Done")
							}
							return err
						},
					},
				},
			},
		},
	}
	app.Run(os.Args)
}
