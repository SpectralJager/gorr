package gorr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/SpectralJager/gorr/bencode"
)

func TestParseTorrent(t *testing.T) {
	input, _ := os.ReadFile("test/torrent.torrent")
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
