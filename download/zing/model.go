package zing

type AlbumResponse struct {
	Err  int
	Msg  string
	Data AlbumResponseData
}

type AlbumResponseData struct {
	Items []AlbumItem
	IsVip bool `json:"is_vip"`
}

type AlbumItem struct {
	Id        string
	Name      string
	Title     string
	Code      string
	Artists   []Artist
	Performer string
	Type      string
	Link      string
	Source    AlbumItemSource
	Order     interface{}
	IsVip     bool `json:"is_vip"`
}

type AlbumItemSource struct {
	Normal string `json:"128"`
	High   string `json:"320"`
}

type Artist struct {
	Name string
	Link string
}

type SongResponse struct {
	Err  int
	Msg  string
	Data AlbumItem
}
