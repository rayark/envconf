package envconf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Load loads config from environment variables into the provided
// pointer-to-struct `out`.
// The names of loaded environment variables are uppercase and all start with the given `prefix`.
//
// Warning:
// 1. Fields without env tag will be ignored.
// 2. Duplicated field tags will be assigned with the same environment variable.
func Load(prefix string, out interface{}, opts ...Option) {
	loader := loader{}
	loader.envStatuses = map[string]*EnvStatus{}
	loader.handleEnvironmentVariables = func(map[string]*EnvStatus) {}
	for _, opt := range opts {
		opt(&loader)
	}

	val := reflect.ValueOf(out).Elem()
	loader.loadStruct(strings.ToUpper(prefix), &val)

	loader.handleEnvironmentVariables(loader.envStatuses)
}

type loader struct {
	handleEnvironmentVariables func(map[string]*EnvStatus)
	envStatuses                map[string]*EnvStatus
}

func (l *loader) useKey(name string) error {
	if _, ok := l.envStatuses[name]; ok {
		return fmt.Errorf("Duplicated key %v", name)
	}
	l.envStatuses[name] = &EnvStatus{}
	return nil
}

type EnvStatus struct {
	Loaded bool
}

func (e *EnvStatus) SetLoaded() {
	e.Loaded = true
}

var stringSliceType = reflect.TypeOf([]string{})

func (l *loader) loadField(name string, out *reflect.Value) {
	kind := out.Kind()
	if kind == reflect.Struct {
		l.loadStruct(name, out)
		return
	}

	if err := l.useKey(name); err != nil {
		panic(err)
	}

	switch out.Kind() {
	case reflect.String:
		l.loadString(name, out)
		return

	case reflect.Int64:
		switch out.Type() {
		case reflect.TypeOf(time.Duration(0)):
			l.loadDuration(name, out)
			return
		}
		fallthrough

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32:
		l.loadInt(name, out)
		return

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		l.loadUint(name, out)
		return

	case reflect.Bool:
		l.loadBool(name, out)
		return

	case reflect.Slice:
		sliceType := out.Type()
		switch sliceType {
		case stringSliceType:
			l.loadStringSlice(name, out)
			return

		default:
			panic(fmt.Errorf("slice type %v on %v is not supported", sliceType, name))
		}

	default:
		panic(fmt.Errorf("field type of %s cannot be recognized by envconf", name))
	}
}

func (l *loader) loadStruct(prefix string, out *reflect.Value) {
	t := out.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := field.Tag.Get("env")
		// ignore "-"
		if tag == "-" {
			continue
		}
		nameAndOpts := strings.Split(tag, ",")
		inline := false
		if len(nameAndOpts) > 1 {
			for _, opt := range nameAndOpts[1:] {
				switch opt {
				case "inline":
					inline = true
				}
			}
		}

		var name string
		if inline {
			if field.Type.Kind() != reflect.Struct {
				panic("Option ,inline needs a struct value field")
			}
			name = prefix
		} else {
			name = nameAndOpts[0]
			if name == "" {
				name = field.Name
			}
			name = prefix + "_" + strings.ToUpper(name)
		}

		fval := out.Field(i)
		l.loadField(name, &fval)
	}
}

func (l *loader) loadString(name string, out *reflect.Value) {
	data, found := syscall.Getenv(name)
	if !found {
		return
	}
	l.envStatuses[name].SetLoaded()

	out.SetString(data)
}

func (l *loader) loadDuration(name string, out *reflect.Value) {
	data, found := syscall.Getenv(name)
	if !found {
		return
	}
	l.envStatuses[name].SetLoaded()

	d, err := time.ParseDuration(data)
	if err != nil {
		panic(fmt.Errorf("field type of %s cannot be parsed into time duration", name))
	}
	out.SetInt(int64(d))
}

func (l *loader) loadInt(name string, out *reflect.Value) {
	data, found := syscall.Getenv(name)
	if !found {
		return
	}
	l.envStatuses[name].SetLoaded()

	d, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		panic(fmt.Errorf("field type of %s cannot be parsed into integer", name))
	}

	out.SetInt(d)
}

func (l *loader) loadUint(name string, out *reflect.Value) {
	data, found := syscall.Getenv(name)
	if !found {
		return
	}
	l.envStatuses[name].SetLoaded()

	d, err := strconv.ParseUint(data, 10, 64)
	if err != nil {
		panic(fmt.Errorf("field type of %s cannot be parsed into unsigned integer", name))
	}

	out.SetUint(d)
}

func (l *loader) loadBool(name string, out *reflect.Value) {
	data, found := syscall.Getenv(name)
	if !found {
		return
	}
	l.envStatuses[name].SetLoaded()

	isTrue := data == "true" || data == "TRUE" || data == "True" || data == "1"
	isFalse := data == "false" || data == "FALSE" || data == "False" || data == "0"

	if !isTrue && !isFalse {
		panic(fmt.Errorf("envvar %s should be a boolean, but cannot be recognized", name))
	}

	out.SetBool(isTrue)
}

func (l *loader) loadStringSlice(name string, out *reflect.Value) {
	data, found := syscall.Getenv(name)
	if !found {
		return
	}
	l.envStatuses[name].SetLoaded()

	strListRaw := strings.Split(data, ",")
	strList := []string{}

	for _, str := range strListRaw {
		v := strings.TrimSpace(str)
		if len(v) > 0 {
			strList = append(strList, v)
		}
	}

	out.Set(reflect.ValueOf(strList))
}
