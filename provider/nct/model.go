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
}
