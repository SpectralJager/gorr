package gorr

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/SpectralJager/gorr/bencode"
)

const PeerID = "-GR0001-000000000000"

type Torrent struct {
	Announce     string     `ben:"announce"`
	AnnounceList [][]string `ben:"announce-list"`
	CreationDate int        `ben:"creation date"`
	Comment      string     `ben:"comment"`
	Encoding     string     `ben:"encoding"`
	Info         Info       `ben:"info"`
}

type Info struct {
	Hash        [20]byte
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

func Open(path string) (Torrent, error) {
	input, _ := os.ReadFile(path)
	document, err := bencode.NewDecoder(bytes.NewReader(input), len(input)).Decode()
	if err != nil {
		return Torrent{}, err
	}
	var torrent Torrent
	err = bencode.Unmarshal(&torrent, document)
	if err != nil {
		return Torrent{}, err
	}
	var buff bytes.Buffer
	err = bencode.NewEncoder(&buff, torrent.Info.Length+500).Encode(document.Get("info"))
	if err != nil {
		return Torrent{}, err
	}
	torrent.Info.Hash = sha1.Sum(buff.Bytes())
	return torrent, nil
}

type Peer struct {
	IP   net.IP
	Port uint16
}

func getPeers(torrent Torrent) ([]Peer, error) {
	req, err := buildRequest(torrent)
	if err != nil {
		return []Peer{}, nil
	}
	peers, err := makeRequest(req)
	if err != nil {
		return []Peer{}, nil
	}
	return peers, nil
}

func buildRequest(torrent Torrent) (*http.Request, error) {
	params := url.Values{
		"info_hash":  []string{string(torrent.Info.Hash[:])},
		"peer_id":    []string{PeerID},
		"uploaded":   []string{strconv.Itoa(0)},
		"downloaded": []string{strconv.Itoa(0)},
		"left":       []string{strconv.Itoa(torrent.Info.Length)},
		"port":       []string{strconv.Itoa(9090)},
		"compact":    []string{strconv.Itoa(1)},
	}
	return http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", torrent.Announce, params.Encode()), nil)
}

func makeRequest(req *http.Request) ([]Peer, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []Peer{}, err
	}
	doc, err := bencode.NewDecoder(resp.Body, int(resp.ContentLength)).Decode()
	if err != nil {
		return []Peer{}, err
	}
	type TrackerResponse struct {
		Interval int    `ben:"interval"`
		Peers    string `ben:"peers"`
	}
	var trackerResponse TrackerResponse
	err = bencode.Unmarshal(&trackerResponse, doc)
	if err != nil {
		return []Peer{}, err
	}
	peers := []Peer{}
	for i := 0; i < len(trackerResponse.Peers); {
		ip := trackerResponse.Peers[i : i+4]
		i += 4
		port := trackerResponse.Peers[i : i+2]
		peers = append(peers, Peer{
			IP:   net.IP(ip),
			Port: binary.BigEndian.Uint16([]byte(port)),
		})
		i += 2
	}
	return peers, nil
}
