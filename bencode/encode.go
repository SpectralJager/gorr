package bencode

import (
	"bufio"
	"fmt"
	"io"
)

type Encoder struct {
	buf *bufio.Writer
}

func NewEncoder(w io.Writer, size int) *Encoder {
	return &Encoder{
		buf: bufio.NewWriterSize(w, size),
	}
}

func (e *Encoder) Encode(val Bencode) error {
	defer e.buf.Flush()
	switch val.typ {
	case INTEGER:
		return e.encodeInteger(val)
	case STRING:
		return e.encodeString(val)
	case LIST:
		return e.encodeList(val)
	case DICTIONARY:
		return e.encodeDictionary(val)
	default:
		return fmt.Errorf("can't encode unsupported bencode type: %d", val.typ)
	}
}

func (e *Encoder) encodeInteger(val Bencode) error {
	_, err := e.buf.WriteString(fmt.Sprintf("i%de", val.integer))
	if err != nil {
		return fmt.Errorf("can't encode integer: %w", err)
	}
	return nil
}

func (e *Encoder) encodeString(val Bencode) error {
	_, err := e.buf.WriteString(fmt.Sprintf("%d:%s", len(val.str), val.str))
	if err != nil {
		return fmt.Errorf("can't encode string: %w", err)
	}
	return nil
}

func (e *Encoder) encodeList(val Bencode) error {
	fmt.Fprint(e.buf, "l")
	for _, item := range val.list {
		err := e.Encode(item)
		if err != nil {
			return fmt.Errorf("can't encode list item: %w", err)
		}
	}
	fmt.Fprint(e.buf, "e")
	return nil
}

func (e *Encoder) encodeDictionary(val Bencode) error {
	fmt.Fprint(e.buf, "d")
	for _, pair := range val.dict.values {
		err := e.encodeString(NewString(pair.Key))
		if err != nil {
			return fmt.Errorf("can't encode list item: %w", err)
		}
		err = e.Encode(pair.Value)
		if err != nil {
			return fmt.Errorf("can't encode list item: %w", err)
		}
	}
	fmt.Fprint(e.buf, "e")
	return nil
}
