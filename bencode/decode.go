package bencode

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Decoder struct {
	buf *bufio.Reader
}

func NewDecoder(r io.Reader, size int) *Decoder {
	return &Decoder{
		buf: bufio.NewReaderSize(r, size),
	}
}

func (d *Decoder) Decode() (Bencode, error) {
	typ, err := d.buf.Peek(1)
	if err != nil {
		return Bencode{}, fmt.Errorf("can't decode type of bencode: %w", err)
	}
	switch typ[0] {
	case 'd':
		return d.decodeDictionary()
	case 'i':
		return d.decodeInteger()
	case 'l':
		return d.decodeList()
	default:
		return d.decodeString()
	}
}

func (d *Decoder) decodeInteger() (Bencode, error) {
	val := Bencode{typ: INTEGER}
	_, err := d.buf.ReadByte()
	if err != nil {
		return Bencode{}, fmt.Errorf("can't decode first byte of integer: %w", err)
	}
	content, err := d.buf.ReadString('e')
	if err != nil {
		return Bencode{}, fmt.Errorf("can't extract integer content: %w", err)
	}
	val.integer, err = strconv.Atoi(content[:len(content)-1])
	if err != nil {
		return Bencode{}, fmt.Errorf("can't decode integer content: %w", err)
	}
	return val, nil
}

func (d *Decoder) decodeString() (Bencode, error) {
	val := Bencode{typ: STRING}
	lenghtStr, err := d.buf.ReadString(':')
	if err != nil {
		return Bencode{}, fmt.Errorf("can't read length of string: %w", err)
	}
	length, err := strconv.Atoi(lenghtStr[:len(lenghtStr)-1])
	if err != nil {
		return Bencode{}, fmt.Errorf("can't decode strings's length: %w", err)
	}
	if length < 0 {
		return Bencode{}, fmt.Errorf("length of string should be >= 0, got: %d", length)
	}
	content := make([]byte, length)
	n, err := d.buf.Read(content)
	if err != nil {
		return Bencode{}, fmt.Errorf("can't extract string content: %w", err)
	}
	if n != length {
		return Bencode{}, fmt.Errorf("can't read whole string")
	}
	val.str = string(content)
	return val, nil
}

func (d *Decoder) decodeList() (Bencode, error) {
	val := Bencode{typ: LIST, list: []Bencode{}}
	_, err := d.buf.ReadByte()
	if err != nil {
		return Bencode{}, fmt.Errorf("can't decode first byte of list: %w", err)
	}
	for {
		if end, err := d.buf.Peek(1); end[0] == 'e' && err == nil {
			break
		}
		item, err := d.Decode()
		if err != nil {
			return Bencode{}, fmt.Errorf("can't decode list item: %w", err)
		}
		val.list = append(val.list, item)
	}
	_, err = d.buf.ReadByte()
	if err != nil {
		return Bencode{}, fmt.Errorf("can't decode last byte of list: %w", err)
	}
	return val, nil
}

func (d *Decoder) decodeDictionary() (Bencode, error) {
	val := Bencode{typ: DICTIONARY, dict: newOrderedMap()}
	_, err := d.buf.ReadByte()
	if err != nil {
		return Bencode{}, fmt.Errorf("can't decode first byte of list: %w", err)
	}
	for {
		if end, err := d.buf.Peek(1); end[0] == 'e' && err == nil {
			break
		}
		key, err := d.decodeString()
		if err != nil {
			return Bencode{}, fmt.Errorf("can't decode dictionary key: %w", err)
		}
		item, err := d.Decode()
		if err != nil {
			return Bencode{}, fmt.Errorf("can't decode dictionary item: %w", err)
		}
		val.dict.set(key.str, item)
	}
	_, err = d.buf.ReadByte()
	if err != nil {
		return Bencode{}, fmt.Errorf("can't decode last byte of list: %w", err)
	}
	return val, nil
}
