package bencode

import (
	"bytes"
	"reflect"
	"testing"
)

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
