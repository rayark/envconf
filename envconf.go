package envconf

import (
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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		loadInteger(name, out)
	case reflect.Struct:
		loadStruct(name, out)
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

func loadInteger(name string, out *reflect.Value) {
	data, found := syscall.Getenv(name)
	if !found {
		return
	}
	d, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
	}

	out.SetInt(d)
}
