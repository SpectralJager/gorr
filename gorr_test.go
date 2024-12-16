package gorr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/SpectralJager/gorr/bencode"
)

func TestGorr(t *testing.T) {
	torrent, err := Open("test/debian.torrent")
	if err != nil {
		t.Fatal(err)
	}
	peers, err := GetPeers(torrent)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.MarshalIndent(peers, "", "  ")
	fmt.Println(string(data))
}

func TestParseTorrent(t *testing.T) {
	input, _ := os.ReadFile("test/miside.torrent")
	document, err := bencode.NewDecoder(bytes.NewReader(input), len(input)).Decode()
	if err != nil {
		t.Fatal(err)
	}
	// var torrent map[string]any
	var torrent Torrent
	err = bencode.Unmarshal(&torrent, document)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.MarshalIndent(torrent, "", "  ")
	fmt.Println(string(data))
}
