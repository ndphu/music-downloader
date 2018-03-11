package provider

import (
	"errors"
	"gopkg.in/urfave/cli.v2"
)

type LoginContext struct {
	UserName string
	Password string
}

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
