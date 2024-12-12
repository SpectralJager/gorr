package bencode

type orderedMap struct {
	values  []Pair
	indexes map[string]int
}

type Pair struct {
	Key   string
	Value Bencode
}

func newOrderedMap() orderedMap {
	return orderedMap{
		indexes: map[string]int{},
		values:  []Pair{},
	}
}

func (m *orderedMap) set(key string, value Bencode) {
	if index, ok := m.indexes[key]; ok {
		m.values[index] = Pair{Key: key, Value: value}
		return
	}
	index := len(m.values)
	m.indexes[key] = index
	m.values = append(m.values, Pair{Key: key, Value: value})
}

func (m *orderedMap) get(key string) Bencode {
	if index, ok := m.indexes[key]; ok {
		return m.values[index].Value
	}
	return Bencode{}
}
