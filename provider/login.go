package provider

import (
	"errors"
	"gopkg.in/urfave/cli.v2"
)

type LoginContext struct {
	UserName string
	Password string
}

var (
	LoginFlags = []cli.Flag{
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
	}
)

func (s *ProviderService) Login(c *cli.Context) error {
	hadProvider := false

	context := &LoginContext{
		UserName: c.String("username"),
		Password: c.String("password"),
	}

	providerName := c.String("name")
	for _, p := range s.providers {
		if p.GetName() == providerName {
			hadProvider = true
			err := p.Login(context)
			if err != nil {
				return err
			}
			break
		}
	}
	if !hadProvider {
		return errors.New("No provider supported: " + providerName)
	}
	return nil
}

func (s *ProviderService) LoginActionHandler(c *cli.Context) error {
	return s.Login(c)
}
