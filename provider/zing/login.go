package zing

import (
	"errors"
	"github.com/ndphu/music-downloader/provider"
)

func (p *ZingProvider) Login(c *provider.LoginContext) error {
	return errors.New("Login for ZING is not supported at this moment!")
}
