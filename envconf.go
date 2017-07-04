package envconf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"syscall"
)

func Load(prefix string, out interface{}) {
	val := reflect.ValueOf(out).Elem()
	loadStruct(strings.ToUpper(prefix), &val)
}

func loadField(name string, out *reflect.Value) {
	switch out.Kind() {
	case reflect.String:
		loadString(name, out)

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		loadInt(name, out)

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		loadUint(name, out)

	case reflect.Bool:
		loadBool(name, out)

	case reflect.Struct:
		loadStruct(name, out)

	default:
		panic(fmt.Errorf("field type of %s cannot be recognized by envconf", name))
	}
}

func loadStruct(prefix string, out *reflect.Value) {
	t := out.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			continue
		}
		name := field.Tag.Get("env")
		if name == "" {
			continue
		}

		name = prefix + "_" + strings.ToUpper(name)
		fval := out.Field(i)
		loadField(name, &fval)
	}
}

func loadString(name string, out *reflect.Value) {
	data, found := syscall.Getenv(name)
	if !found {
		return
	}

	out.SetString(data)
}

func loadInt(name string, out *reflect.Value) {
	data, found := syscall.Getenv(name)
	if !found {
		return
	}

	d, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		panic(fmt.Errorf("field type of %s cannot be parsed into integer", name))
	}

	out.SetInt(d)
}

func loadUint(name string, out *reflect.Value) {
	data, found := syscall.Getenv(name)
	if !found {
		return
	}

	d, err := strconv.ParseUint(data, 10, 64)
	if err != nil {
	}

	out.SetUint(d)
}

func loadBool(name string, out *reflect.Value) {
	data, found := syscall.Getenv(name)
	if !found {
		return
	}

	isTrue := data == "true" || data == "TRUE" || data == "True" || data == "1"
	isFalse := data == "false" || data == "FALSE" || data == "False" || data == "0"

	if !isTrue && !isFalse {
		panic(fmt.Errorf("envvar %s should be a boolean, but cannot be recognized", name))
	}

	out.SetBool(isTrue)
}
