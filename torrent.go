package gorr

type Torrent struct {
	Announce     string     `ben:"announce"`
	AnnounceList [][]string `ben:"announce-list"`
	CreationDate int        `ben:"creation date"`
	Comment      string     `ben:"comment"`
	Encoding     string     `ben:"encoding"`
	Info         Info       `ben:"info"`
}

type Info struct {
	PieceLength int    `ben:"piece length"`
	Pieces      string `ben:"pieces"`
	Private     bool   `ben:"private"`
	Path        string `ben:"name"`
	Length      int    `ben:"length"`
	Files       []struct {
		Length int      `ben:"length"`
		Path   []string `ben:"path"`
	} `ben:"files"`
}
