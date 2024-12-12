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
	torrentFile := NewTorrent(document)
	data, _ := json.MarshalIndent(torrentFile, "", "  ")
	fmt.Println(string(data))
}
