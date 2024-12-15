package bencode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
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

func TestDecodeInteger(t *testing.T) {
	testCases := []struct {
		input  string
		expect Bencode
	}{
		{
			input:  "i13e",
			expect: Bencode{typ: INTEGER, integer: 13},
		},
		{
			input:  "i-13e",
			expect: Bencode{typ: INTEGER, integer: -13},
		},
		{
			input:  "i0e",
			expect: Bencode{typ: INTEGER, integer: 0},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			val, err := NewDecoder(strings.NewReader(tC.input), 0).decodeInteger()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tC.expect, val) {
				t.Fail()
			}
		})
	}
}

func TestDecodeString(t *testing.T) {
	testCases := []struct {
		input  string
		expect Bencode
	}{
		{
			input:  "5:hello",
			expect: Bencode{typ: STRING, str: "hello"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			val, err := NewDecoder(strings.NewReader(tC.input), 0).decodeString()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tC.expect, val) {
				t.Fail()
			}
		})
	}
}

func TestDecodeList(t *testing.T) {
	testCases := []struct {
		input  string
		expect Bencode
	}{
		{
			input: "l5:hello5:worlde",
			expect: Bencode{typ: LIST, list: []Bencode{
				{typ: STRING, str: "hello"},
				{typ: STRING, str: "world"},
			}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			val, err := NewDecoder(strings.NewReader(tC.input), 0).decodeList()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tC.expect, val) {
				t.Fail()
			}
		})
	}
}

func TestDecodeDictionary(t *testing.T) {
	testCases := []struct {
		input  string
		expect Bencode
	}{
		{
			input: "d4:name6:Daniele",
			expect: NewDictionary(
				Pair{Key: "name", Value: NewString("Daniel")},
			),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			val, err := NewDecoder(strings.NewReader(tC.input), 0).decodeDictionary()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tC.expect, val) {
				t.Fail()
			}
		})
	}
}

func TestEncodeInteger(t *testing.T) {
	testCases := []struct {
		input  Bencode
		expect string
	}{
		{
			input:  NewInteger(10),
			expect: "i10e",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.expect, func(t *testing.T) {
			data := bytes.NewBuffer([]byte{})
			err := NewEncoder(data, 0).Encode(tC.input)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tC.expect, string(data.Bytes())) {
				t.Fail()
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	input := `d4:numsd1:1i1e1:2i2eee`
	document, err := NewDecoder(bytes.NewBufferString(input), len(input)).Decode()
	if err != nil {
		t.Fatal(err)
	}
	var val struct {
		Lol  string
		Nums struct {
			First  int `ben:"1"`
			Second int `ben:"2"`
		} `ben:"nums"`
	}
	err = Unmarshal(&val, document)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(val)
	bs, _ := json.Marshal(val)
	fmt.Println(string(bs))
}
