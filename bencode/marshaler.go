package bencode

import (
	"fmt"
	"reflect"
)

func Unmarshal(val any, doc Bencode) error {
	refVal := reflect.ValueOf(val)
	if refVal.Kind() != reflect.Ptr {
		return fmt.Errorf("can't unmarshal into non pointer")
	}
	return unmarshal(refVal.Elem(), doc)
}

func unmarshal(val reflect.Value, doc Bencode) error {
	switch val.Kind() {
	case reflect.Int:
		if doc.Type() != INTEGER {
			return fmt.Errorf("can't unmarshal '%d' bencode to integer", doc.Type())
		}
		val.SetInt(int64(doc.Integer()))
	case reflect.Bool:
		if doc.Type() != INTEGER {
			return fmt.Errorf("can't unmarshal '%d' bencode to bool", doc.Type())
		}
		if doc.Integer() != 0 {
			val.SetBool(true)
		} else {
			val.SetBool(false)
		}
	case reflect.String:
		if doc.Type() != STRING {
			return fmt.Errorf("can't unmarshal '%d' bencode to string", doc.Type())
		}
		val.SetString(doc.Str())
	case reflect.Slice:
		if doc.Type() != LIST {
			return fmt.Errorf("can't unmarshal '%d' bencode to slice", doc.Type())
		}
		val.Set(reflect.MakeSlice(val.Type(), doc.Len(), doc.Len()))
		for i := 0; i < doc.Len(); i++ {
			err := unmarshal(val.Index(i), doc.Item(i))
			if err != nil {
				return fmt.Errorf("can't unmarshal %dth slice item: %w", i, err)
			}
		}
	case reflect.Map:
		typ := val.Type()
		if val.IsNil() {
			val.Set(reflect.MakeMap(typ))
		}
		if typ.Key().Kind() != reflect.String {
			return fmt.Errorf("can't unmarshal to map: type of key should be string")
		}
		if typ.Elem().Kind() != reflect.Interface && typ.Elem().NumMethod() != 0 {
			return fmt.Errorf("can't unmarshal to map: type of value should be interface{}")
		}
		if doc.Type() != DICTIONARY {
			return fmt.Errorf("can't unmarshal '%d' bencode to map", doc.Type())
		}
		keys := doc.Keys()
		for _, key := range keys {
			elem := reflect.New(typ.Elem())
			err := unmarshal(elem.Elem(), doc.Get(key))
			if err != nil {
				return fmt.Errorf("can't unmarshal '%s' map item: %w", key, err)
			}
			val.SetMapIndex(reflect.ValueOf(key), elem)
		}
	case reflect.Interface:
		typ := val.Type()
		if typ.NumMethod() != 0 {
			return fmt.Errorf("can't unmarshal to %T, expect interface{}", typ.String())
		}

		switch doc.Type() {
		case INTEGER:
			val.Set(reflect.ValueOf(doc.Integer()).Convert(typ))
		case STRING:
			val.Set(reflect.ValueOf(doc.Str()).Convert(typ))
		case LIST:
			list := reflect.MakeSlice(reflect.TypeOf([]any{}), doc.Len(), doc.Len())
			for i := 0; i < doc.Len(); i++ {
				err := unmarshal(list.Index(i), doc.Item(i))
				if err != nil {
					return fmt.Errorf("can't unmarshal to []any: %w", err)
				}
			}
			val.Set(list.Convert(typ))
		case DICTIONARY:
			dict := reflect.MakeMap(reflect.TypeOf(map[string]any{}))
			err := unmarshal(dict, doc)
			if err != nil {
				return fmt.Errorf("can't unmarshal to map[string]any: %w", err)
			}
			val.Set(dict.Convert(typ))
		default:
			return fmt.Errorf("can't unmarshal '%d' to %s", doc.Type(), typ.String())
		}
	case reflect.Struct:
		if doc.Type() != DICTIONARY {
			return fmt.Errorf("can't unmarshal '%d' bencode to struct", doc.Type())
		}
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			fieldType := typ.Field(i)
			tagName := fieldType.Tag.Get("ben")
			if len(tagName) == 0 {
				continue
			}
			item := doc.Get(tagName)
			if item.Type() == ILLEGAL {
				continue
			}
			err := unmarshal(val.Field(i), item)
			if err != nil {
				return fmt.Errorf("can't unmarshal to struct field: %w", err)
			}
		}
	default:
		return fmt.Errorf("can't unmarshal to %s: unsupported", val.Type().String())
	}
	return nil
}
