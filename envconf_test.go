package envconf

import (
	"os"
	"reflect"
	"testing"
	"time"
)

type Config struct {
	Mongo          MongoConfig `env:"mongo"`
	AppIDList      []string    `env:"app_id_list"`
	Replicas       uint        `env:"replicas"`
	EmbeddedConfig `env:",inline"`

	Duration time.Duration `env:"duration"`
	I64      int64         `env:"i64"`

	unexported       string
	UnTagged         float64       `env:"-"` // unsupported types must be ignored explicitly
	UnLoadedInt      int           `env:"unloadedint"`
	UnLoadedUint     uint          `env:"unloadeduint"`
	UnLoadedBool     bool          `env:"unloadedbool"`
	UnLoadedStr      string        `env:"unloadedstring"`
	UnLoadedStrSlice []string      `env:"unloadedstrslice"`
	UnLoadedDuration time.Duration `env:"unloadedduration"`
}

type MongoConfig struct {
	Nodes      string `env:"nodes"`
	Database   string `env:"db"`
	ReplicaSet string `env:"replicaset"`
	Port       int    `env:"port"`
	Debug      bool   `env:"debug"`
}

type EmbeddedConfig struct {
	StringInEmbeddedStructure string `env:"string_in_embedded_structure"`
	IntInEmbeddedStructure    int    `env:"int_in_embedded_structure"`
}

func assertEqual(t testing.TB, fieldName string, expected, received interface{}) {
	if !reflect.DeepEqual(expected, received) {
		t.Helper()
		t.Errorf("%v is not loaded correctly.\nExpecting: %v %v\nReceived:  %v %v",
			fieldName,
			reflect.TypeOf(expected), expected,
			reflect.TypeOf(received), received,
		)
	}
}
func assertPanic(t testing.TB) {
	if r := recover(); r == nil {
		t.Helper()
		t.Errorf("this test should panic")
	}
}

func TestLoad(t *testing.T) {
	os.Setenv("TEST_MONGO_NODES", "www.example.com")
	os.Setenv("TEST_MONGO_PORT", "332")
	os.Setenv("TEST_MONGO_DEBUG", "false")
	os.Setenv("TEST_APP_ID_LIST", " aa, bb ,cc ,dd")
	os.Setenv("TEST_REPLICAS", "3")
	os.Setenv("TEST_STRING_IN_EMBEDDED_STRUCTURE", "a-z")
	os.Setenv("TEST_INT_IN_EMBEDDED_STRUCTURE", "19")
	os.Setenv("TEST_DURATION", "10m")
	os.Setenv("TEST_I64", "600")

	initConfig := Config{
		unexported:       "unexported string",
		UnTagged:         98.76,
		UnLoadedInt:      -5,
		UnLoadedUint:     9,
		UnLoadedBool:     true,
		UnLoadedStr:      "unloaded string",
		UnLoadedStrSlice: []string{"some", "random", "str"},
	}

	config := initConfig
	Load("TEST", &config)

	assertEqual(t, "Mongo.Port", 332, config.Mongo.Port)
	assertEqual(t, "Mongo.Nodes", "www.example.com", config.Mongo.Nodes)
	assertEqual(t, "Mongo.Database", "", config.Mongo.Database)
	assertEqual(t, "Mongo.Debug", false, config.Mongo.Debug)
	assertEqual(t, "Replicas", uint(3), config.Replicas)
	assertEqual(t, "AppIDList", []string{"aa", "bb", "cc", "dd"}, config.AppIDList)
	assertEqual(t, "StringInEmbeddedStrcuture", "a-z", config.StringInEmbeddedStructure)
	assertEqual(t, "IntInEmbeddedStrcuture", 19, config.IntInEmbeddedStructure)
	assertEqual(t, "Duration", time.Minute*10, config.Duration)
	assertEqual(t, "int64", int64(600), config.I64)

	assertEqual(t, "unexported", initConfig.unexported, config.unexported)
	assertEqual(t, "UnTagged", initConfig.UnTagged, config.UnTagged)
	assertEqual(t, "UnLoadedInt", initConfig.UnLoadedInt, config.UnLoadedInt)
	assertEqual(t, "UnLoadedUint", initConfig.UnLoadedUint, config.UnLoadedUint)
	assertEqual(t, "UnLoadedBool", initConfig.UnLoadedBool, config.UnLoadedBool)
	assertEqual(t, "UnLoadedStr", initConfig.UnLoadedStr, config.UnLoadedStr)
	assertEqual(t, "UnLoadedStrSlice", initConfig.UnLoadedStrSlice, config.UnLoadedStrSlice)
}

func TestTaggedUnsupportedTypeShouldPanic(t *testing.T) {
	defer assertPanic(t)

	type Invalid struct {
		Unsupported float32 `env:"unsupported"`
	}

	os.Setenv("FAIL_UNSUPPORTED", "55.66")
	var invalid Invalid
	Load("FAIL", &invalid)
}

func TestUnsupportedSliceShouldPanic(t *testing.T) {
	defer assertPanic(t)

	type Invalid struct {
		Unsupported []map[string]string
	}

	os.Setenv("FAIL_UNSUPPORTED", "[]")
	var invalid Invalid
	Load("FAIL", &invalid)
}

func TestInvalidIntShouldPanic(t *testing.T) {
	defer assertPanic(t)

	type InvalidInt struct {
		InvalidInt int `env:"invalidint"`
	}

	os.Setenv("FAIL_INVALIDINT", "not a int")
	var inv InvalidInt
	Load("FAIL", &inv)
}

func TestInvalidDurationShouldPanic(t *testing.T) {
	defer assertPanic(t)

	type InvalidDuration struct {
		InvalidDuration time.Duration `env:"invalidduration"`
	}

	os.Setenv("FAIL_INVALIDDURATION", "not a duration")
	var inv InvalidDuration
	Load("FAIL", &inv)
}

func TestInvalidUintShouldPanic(t *testing.T) {
	defer assertPanic(t)

	type InvalidUint struct {
		InvalidUint uint `env:"invaliduint"`
	}

	os.Setenv("FAIL_INVALIDUINT", "-2")
	var inv InvalidUint
	Load("FAIL", &inv)
}

func TestInvalidBoolShouldPanic(t *testing.T) {
	defer assertPanic(t)

	type InvalidBool struct {
		InvalidBool bool `env:"invalidbool"`
	}

	os.Setenv("FAIL_INVALIDBOOL", "ture")
	var inv InvalidBool
	Load("FAIL", &inv)
}

func TestInlineWithoutComma(t *testing.T) {
	os.Setenv("TEST_STRING_IN_EMBEDDED_STRUCTURE", "a-z")
	os.Setenv("TEST_INT_IN_EMBEDDED_STRUCTURE", "19")

	// no comma
	config := struct {
		EmbeddedConfig `env:"inline"`
	}{}
	Load("TEST", &config)

	assertEqual(t, "StringInEmbeddedStructure", "", config.StringInEmbeddedStructure)
	assertEqual(t, "IntInEmbeddedStructure", 0, config.IntInEmbeddedStructure)
}

func TestInlineWithTagName(t *testing.T) {
	os.Setenv("TEST_STR", "outer")
	os.Setenv("TEST_STRING_IN_EMBEDDED_STRUCTURE", "a-z")
	os.Setenv("TEST_INT_IN_EMBEDDED_STRUCTURE", "19")
	os.Setenv("TEST_STRING_IN_STRUCTURE", "in struct")
	os.Setenv("TEST_INT_IN_STRUCTURE", "1239")

	// inline option ignores tag name
	type structure struct {
		StringInStructure string `env:"string_in_structure"`
		IntInStructure    int    `env:"int_in_structure"`
	}
	config := struct {
		Str            string `env:"str"`
		EmbeddedConfig `env:"str,inline"`
		Struct         structure `env:"str,inline"`
	}{}
	Load("TEST", &config)

	assertEqual(t, "Str", "outer", config.Str)
	assertEqual(t, "StringInEmbeddedStructure", "a-z", config.StringInEmbeddedStructure)
	assertEqual(t, "IntInEmbeddedStructure", 19, config.IntInEmbeddedStructure)
	assertEqual(t, "StringInStructure", "in struct", config.Struct.StringInStructure)
	assertEqual(t, "IntInStructure", 1239, config.Struct.IntInStructure)
}

func TestInlineWithoutStruct(t *testing.T) {
	defer assertPanic(t)

	config := struct {
		Str string `env:",inline"`
	}{}
	Load("TEST", &config)
}

func TestEmbeddedStructureWithoutInline(t *testing.T) {
	os.Setenv("TEST_EMBEDDED_STRING_IN_EMBEDDED_STRUCTURE", "a-z")
	os.Setenv("TEST_EMBEDDED_INT_IN_EMBEDDED_STRUCTURE", "19")

	configWithTag := struct {
		EmbeddedConfig `env:"embedded"`
	}{}
	Load("TEST", &configWithTag)
	assertEqual(t, "StringInEmbeddedStructure", "a-z", configWithTag.StringInEmbeddedStructure)
	assertEqual(t, "IntInEmbeddedStructure", 19, configWithTag.IntInEmbeddedStructure)
}

func TestIgnoringSpecificTag(t *testing.T) {
	// invalid
	os.Setenv("TEST_-", "panics")
	// ignored
	os.Setenv("TEST_IGNOREDSTRING", "ignored")
	os.Setenv("TEST_IGNOREDINT", "123")
	os.Setenv("TEST_IGNOREDSTRUCT_IGNOREDSTRINGINSTRUCT", "ignored")
	os.Setenv("TEST_IGNOREDSTRUCT_IGNOREDINTINSTRUCT", "12345")
	os.Setenv("TEST_EMBEDDEDCONFIG_STRING_IN_EMBEDDED_STRUCTURE", "a-z")
	os.Setenv("TEST_EMBEDDEDCONFIG_INT_IN_EMBEDDED_STRUCTURE", "19")

	config := struct {
		IgnoredString string `env:"-"`
		IgnoredInt    int    `env:"-"`
		IgnoredStruct struct {
			IgnoredStringInStruct string
			IgnoredIntInStruct    int
		} `env:"-"`
		EmbeddedConfig `env:"-"`
	}{}
	Load("TEST", &config)

	assertEqual(t, "IgnoredString", "", config.IgnoredString)
	assertEqual(t, "IgnoredInt", 0, config.IgnoredInt)
	assertEqual(t, "IgnoredStringInStruct", "", config.IgnoredStruct.IgnoredStringInStruct)
	assertEqual(t, "IgnoredIntInStruct", 0, config.IgnoredStruct.IgnoredIntInStruct)
	assertEqual(t, "StringInEmbeddedStructure", "", config.StringInEmbeddedStructure)
	assertEqual(t, "IntInEmbeddedStructure", 0, config.IntInEmbeddedStructure)
}

func TestNoTag(t *testing.T) {
	os.Setenv("TEST_AUTONAMEDSTRING", "auto named")
	os.Setenv("TEST_AUTONAMEDINT", "321")
	os.Setenv("TEST_AUTONAMEDSTRUCT_AUTONAMEDSTRINGINSTRUCT", "auto named in struct")
	os.Setenv("TEST_AUTONAMEDSTRUCT_AUTONAMEDINTINSTRUCT", "54321")
	os.Setenv("TEST_EMBEDDEDCONFIG_STRING_IN_EMBEDDED_STRUCTURE", "a-z")
	os.Setenv("TEST_EMBEDDEDCONFIG_INT_IN_EMBEDDED_STRUCTURE", "19")

	config := struct {
		AutoNamedString string
		AutoNamedInt    int
		AutoNamedStruct struct {
			AutoNamedStringInStruct string
			AutoNamedIntInStruct    int
		}
		EmbeddedConfig
	}{}
	Load("TEST", &config)

	assertEqual(t, "AutoNamedString", "auto named", config.AutoNamedString)
	assertEqual(t, "AutoNamedInt", 321, config.AutoNamedInt)
	assertEqual(t, "AutoNamedStringInStruct", "auto named in struct", config.AutoNamedStruct.AutoNamedStringInStruct)
	assertEqual(t, "AutoNamedIntInStruct", 54321, config.AutoNamedStruct.AutoNamedIntInStruct)
	assertEqual(t, "StringInEmbeddedStructure", "a-z", config.StringInEmbeddedStructure)
	assertEqual(t, "IntInEmbeddedStructure", 19, config.IntInEmbeddedStructure)
}

func TestDuplicatedKeys(t *testing.T) {
	defer assertPanic(t)
	os.Setenv("TEST_DUPLICATED_STRING", "duplication")

	config := struct {
		Str    string `env:"duplicated_string"`
		Struct struct {
			Str string `env:"string"`
		} `env:"duplicated"`
	}{}
	Load("TEST", &config)
}
func TestDuplicatedKeysBetweenInlineStructs(t *testing.T) {
	defer assertPanic(t)
	os.Setenv("TEST_DUPLICATED_STRING", "duplication")

	type em1 struct {
		Str1 string `env:"duplicated_string"`
	}
	type em2 struct {
		Str2 string `env:"duplicated_string"`
	}
	config := struct {
		em1 `env:",inline"`
		em2 `env:",inline"`
	}{}
	Load("TEST", &config)
}
func TestDuplicatedKeysBetweenStructs(t *testing.T) {
	defer assertPanic(t)
	os.Setenv("TEST_DUPLICATED_STRING", "duplication")

	type em1 struct {
		Str1 string `env:"duplicated_string"`
	}
	type em2 struct {
		Str2 string `env:"string"`
	}
	config := struct {
		EM1 em1 `env:"em"`
		EM2 em2 `env:"em_duplicated"`
	}{}
	Load("TEST", &config)
}

func TestLogger(t *testing.T) {
	os.Setenv("TEST_INTEGER", "-3")
	os.Setenv("TEST_UNSIGNED_INTEGER", "3")

	result := map[string]*EnvStatus{}
	config := struct {
		String  string `env:"string"`
		Integer int    `env:"integer"`
		Struct  struct {
			Unsigned uint     `env:"unsigned_integer"`
			Bool     bool     `env:"bool"`
			StrSlice []string `env:"string_slice"`
		} `env:",inline"`
	}{}
	Load("TEST", &config, CustomHandleEnvVarsOption(func(status map[string]*EnvStatus) {
		result = status
	}))

	expected := map[string]*EnvStatus{
		"TEST_STRING":           {false},
		"TEST_INTEGER":          {true},
		"TEST_UNSIGNED_INTEGER": {true},
		"TEST_BOOL":             {false},
		"TEST_STRING_SLICE":     {false},
	}
	assertEqual(t, "", expected, result)
}
