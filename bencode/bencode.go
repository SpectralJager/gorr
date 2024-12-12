package bencode

const (
	ILLEGAL byte = iota
	STRING
	INTEGER
	DICTIONARY
	LIST
)

type Bencode struct {
	typ     byte
	integer int
	str     string
	list    []Bencode
	dict    orderedMap
}

func NewInteger(val int) Bencode {
	return Bencode{
		typ:     INTEGER,
		integer: val,
	}
}

func NewString(val string) Bencode {
	return Bencode{
		typ: STRING,
		str: val,
	}
}

func NewList(vals ...Bencode) Bencode {
	return Bencode{
		typ:  LIST,
		list: vals,
	}
}

func NewDictionary(pairs ...Pair) Bencode {
	ordMap := newOrderedMap()
	for _, val := range pairs {
		ordMap.set(val.Key, val.Value)
	}
	return Bencode{
		typ:  DICTIONARY,
		dict: ordMap,
	}
}

func (val Bencode) Type() byte {
	return val.typ
}

func (val Bencode) Integer() int {
	if val.typ == INTEGER {
		return val.integer
	}
	return 0
}

func (val Bencode) Str() string {
	if val.typ == STRING {
		return val.str
	}
	return ""
}

func (val Bencode) Len() int {
	switch val.typ {
	case STRING:
		return len(val.str)
	case LIST:
		return len(val.list)
	default:
		return 0
	}
}

func (val Bencode) Item(index int) Bencode {
	if val.typ != LIST {
		return Bencode{}
	}
	if index < 0 || len(val.list) <= index {
		return Bencode{}
	}
	return val.list[index]
}

func (val Bencode) Keys() []string {
	if val.typ != DICTIONARY {
		return []string{}
	}
	keys := []string{}
	for key := range val.dict.indexes {
		keys = append(keys, key)
	}
	return keys
}

func (val Bencode) Get(key string) Bencode {
	if val.typ != DICTIONARY {
		return Bencode{}
	}
	return val.dict.get(key)
}
