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
	// open torrent file
	torrent, err := Open("test/debian.torrent")
	if err != nil {
		t.Fatal(err)
	}
	// get torrent peers
	peers, err := getPeers(torrent)
	if err != nil {
		t.Fatal(err)
	}
	// download pieces of torrent
	file, err := os.OpenFile("test/"+torrent.Info.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	handshake := bytes.Buffer{}
	handshake.WriteByte(19)
	handshake.WriteString("BitTorrent protocol")
	handshake.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0})
	handshake.Write(torrent.Info.Hash[:])
	handshake.WriteString(PeerID)
	for _, peer := range peers {
		_ = peer
	}
	_ = file

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
