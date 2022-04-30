package feeder

import (
	"fmt"
	"github.com/golobby/cast"
	"reflect"
	"unsafe"
)

// Default is a feeder.
// It feeds using default value if tag exists.
type Default struct{}

// Feed set default values into the given struct.
// `default:"value"` is a tag example.
func (f Default) Feed(structure interface{}) error {
	inputType := reflect.TypeOf(structure)
	if inputType != nil {
		if inputType.Kind() == reflect.Ptr {
			if inputType.Elem().Kind() == reflect.Struct {
				return fillStruct(reflect.ValueOf(structure).Elem())
			}
		}
	}

	return nil
}

// fillStruct sets a reflected struct fields with the default value.
func fillStruct(s reflect.Value) error {
	for i := 0; i < s.NumField(); i++ {
		if t, exist := s.Type().Field(i).Tag.Lookup("default"); exist {
			v, err := cast.FromType(t, s.Type().Field(i).Type)
			if err != nil {
				return fmt.Errorf("default: cannot set `%v` field; err: %v", s.Type().Field(i).Name, err)
			}

			ptr := reflect.NewAt(s.Field(i).Type(), unsafe.Pointer(s.Field(i).UnsafeAddr())).Elem()
			ptr.Set(reflect.ValueOf(v))
		} else if s.Type().Field(i).Type.Kind() == reflect.Struct {
			if err := fillStruct(s.Field(i)); err != nil {
				return err
			}
		} else if s.Type().Field(i).Type.Kind() == reflect.Ptr {
			if s.Field(i).IsZero() == false && s.Field(i).Elem().Type().Kind() == reflect.Struct {
				if err := fillStruct(s.Field(i).Elem()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
