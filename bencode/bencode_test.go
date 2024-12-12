package bencode

import (
	"bytes"
	"os"
	"testing"
)

func TestDecodeDocument(t *testing.T) {
	input, err := os.ReadFile("../test/torrent.torrent")
	if err != nil {
		t.Fatal(err)
	}
	document, err := NewDecoder(bytes.NewReader(input), len(input)).Decode()
	if err != nil {
		t.Fatal(err)
	}
	data := bytes.NewBuffer([]byte{})
	err = NewEncoder(data, len(input)).Encode(document)
	if err != nil {
		t.Fatal(err)
	}
	os.WriteFile("../test/test_res.torrent", data.Bytes(), 0666)
}
