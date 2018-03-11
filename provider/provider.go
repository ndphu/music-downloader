package provider

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
