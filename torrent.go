package gorr

import (
	"time"

	"github.com/SpectralJager/gorr/bencode"
)

type Torrent struct {
	Announce     string
	AnnounceList [][]string
	CreationDate time.Time
	Comment      string
	Encoding     string
	Info         Info
}

func NewTorrent(doc bencode.Bencode) Torrent {
	var torrent Torrent
	torrent.Announce = doc.Get("announce").Str()
	announceList := doc.Get("announce-list")
	for i := 0; i < announceList.Len(); i++ {
		item := announceList.Item(i)
		items := []string{}
		for i := 0; i < item.Len(); i++ {
			items = append(items, item.Item(i).Str())
		}
		torrent.AnnounceList = append(torrent.AnnounceList, items)
	}
	torrent.CreationDate = time.Unix(int64(doc.Get("creation date").Integer()), 0)
	torrent.Comment = doc.Get("comment").Str()
	torrent.Encoding = doc.Get("encoding").Str()
	torrent.Info = NewInfo(doc.Get("info"))
	return torrent
}

type Info struct {
	PieceLength int
	Pieces      []byte
	Private     bool
	Multiple    InfoMultipleFiles
	Name        string
	Length      int
	MD5sum      []byte
}

func NewInfo(doc bencode.Bencode) Info {
	var info Info
	info.PieceLength = doc.Get("piece length").Integer()
	info.Pieces = []byte(doc.Get("pieces").Str())
	priv := doc.Get("private").Integer()
	if priv == 0 {
		info.Private = false
	} else {
		info.Private = true
	}
	if ml := doc.Get("files"); ml.Type() != bencode.ILLEGAL {
		info.Multiple = NewInfoMultipleFiles(doc)
	} else {
		info.Name = doc.Get("name").Str()
		info.Length = doc.Get("length").Integer()
		info.MD5sum = []byte(doc.Get("md5sum").Str())
	}
	return info
}

type InfoMultipleFiles struct {
	Directory string
	Files     []InfoFile
}

func NewInfoMultipleFiles(doc bencode.Bencode) InfoMultipleFiles {
	var info InfoMultipleFiles
	info.Directory = doc.Get("name").Str()
	files := doc.Get("files")
	for i := 0; i < files.Len(); i++ {
		info.Files = append(info.Files, NewInfoFile(files.Item(i)))
	}
	return info
}

type InfoFile struct {
	Path   []string
	Length int
	MD5sum []byte
}

func NewInfoFile(doc bencode.Bencode) InfoFile {
	var info InfoFile
	info.Length = doc.Get("length").Integer()
	info.MD5sum = []byte(doc.Get("md5sum").Str())
	paths := doc.Get("path")
	for i := 0; i < paths.Len(); i++ {
		info.Path = append(info.Path, paths.Item(i).Str())
	}
	return info
}
