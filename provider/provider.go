package provider

import (
	"fmt"
	"gopkg.in/urfave/cli.v2"
)

type Provider interface {
	GetName() string
	Login(*LoginContext) error
	GetSupportedSites() []string
	IsSiteSupported(string) bool
	Download(*DownloadContext) error
}

type ProviderService struct {
	providers []Provider
}

func NewProviderService(_providers []Provider) *ProviderService {
	return &ProviderService{
		providers: _providers,
	}
}

func (s *ProviderService) GetProviders() []Provider {
	return s.providers
}

func (s *ProviderService) ListProviderHandler(c *cli.Context) error {
	fmt.Println("Supported providers: ")
	for _, p := range s.GetProviders() {
		fmt.Println("  - " + p.GetName())
	}
	return nil
}
