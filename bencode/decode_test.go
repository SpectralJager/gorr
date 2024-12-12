package bencode

import (
	"reflect"
	"strings"
	"testing"
)

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
