package nct

type Tracklist struct {
	PageTitle string
	Type      string  `xml:"type"`
	Tracks    []Track `xml:"track"`
}

type Track struct {
	Title      string `xml:"title"`
	Creator    string `xml:"creator"`
	Key        string `xml:"key"`
	Location   string `xml:"location"`
	LocationHQ string `xml:"locationHQ"`
	HasHQ      string `xml:"hasHQ"`
	Info       string `xml:"info"`
}

type LosslessResponse struct {
	StatusReadMode bool              `json:"STATUS_READ_MODE"`
	Data           map[string]string `json:"data"`
	ErrorCode      int               `json:"error_code"`
	ErrorMessage   string            `json:"error_message`
}
