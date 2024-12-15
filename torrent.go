package gorr

type Torrent struct {
	Announce     string     `ben:"announce"`
	AnnounceList [][]string `ben:"announce-list"`
	CreationDate int        `ben:"creation date"`
	Comment      string     `ben:"comment"`
	Encoding     string     `ben:"encoding"`
}
