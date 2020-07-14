package envconf

import (
	"os"
	"reflect"
	"testing"
)

type Config struct {
	Mongo          MongoConfig `env:"mongo"`
	AppIDList      []string    `env:"app_id_list"`
	Replicas       uint        `env:"replicas"`
	EmbeddedConfig `env:",inline"`

	unexported       string
	UnTagged         float64  // unsupported types cannot be tagged with env
	UnLoadedInt      int      `env:"unloadedint"`
	UnLoadedUint     uint     `env:"unloadeduint"`
	UnLoadedBool     bool     `env:"unloadedbool"`
	UnLoadedStr      string   `env:"unloadedstring"`
	UnLoadedStrSlice []string `env:"unloadedstrslice"`
}

type MongoConfig struct {
	Nodes      string `env:"nodes"`
	Database   string `env:"db"`
	ReplicaSet string `env:"replicaset"`
	Port       int    `env:"port"`
	Debug      bool   `env:"debug"`
}

type EmbeddedConfig struct {
	StringInEmbeddedStrcuture string `env:"string_in_embedded_structure"`
}

func assert(t testing.TB, field string, expected, received interface{}) {
	if !reflect.DeepEqual(expected, received) {
		t.Errorf("%v is not loaded correctly.\nExpecting: %v %v\nReceived:  %v %v",
			field,
			reflect.TypeOf(expected), expected,
			reflect.TypeOf(received), received,
		)
	}
}

func TestLoad(t *testing.T) {
	os.Setenv("TEST_MONGO_NODES", "www.example.com")
	os.Setenv("TEST_MONGO_PORT", "332")
	os.Setenv("TEST_MONGO_DEBUG", "false")
	os.Setenv("TEST_APP_ID_LIST", " aa, bb ,cc ,dd")
	os.Setenv("TEST_REPLICAS", "3")
	os.Setenv("TEST_STRING_IN_EMBEDDED_STRUCTURE", "a-z")

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

	assert(t, "Mongo.Port", 332, config.Mongo.Port)
	assert(t, "Mongo.Nodes", "www.example.com", config.Mongo.Nodes)
	assert(t, "Mongo.Database", "", config.Mongo.Database)
	assert(t, "Mongo.Debug", false, config.Mongo.Debug)
	assert(t, "Replicas", uint(3), config.Replicas)
	assert(t, "StringInEmbeddedStrcuture", "a-z", config.StringInEmbeddedStrcuture)
	assert(t, "AppIDList", []string{"aa", "bb", "cc", "dd"}, config.AppIDList)

	assert(t, "unexported", initConfig.unexported, config.unexported)
	assert(t, "UnTagged", initConfig.UnTagged, config.UnTagged)
	assert(t, "UnLoadedInt", initConfig.UnLoadedInt, config.UnLoadedInt)
	assert(t, "UnLoadedUint", initConfig.UnLoadedUint, config.UnLoadedUint)
	assert(t, "UnLoadedBool", initConfig.UnLoadedBool, config.UnLoadedBool)
	assert(t, "UnLoadedStr", initConfig.UnLoadedStr, config.UnLoadedStr)
	assert(t, "UnLoadedStrSlice", initConfig.UnLoadedStrSlice, config.UnLoadedStrSlice)
}

func TestTaggedUnsupportedTypeShouldPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("this test should panic")
		}
	}()

	type Invalid struct {
		Unsupported float32 `env:"unsupported"`
	}

	os.Setenv("FAIL_UNSUPPORTED", "55.66")
	var invalid Invalid
	Load("FAIL", &invalid)
}

func TestInvalidIntShouldPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("this test should panic")
		}
	}()

	type InvalidInt struct {
		InvalidInt int `env:"invalidint"`
	}

	os.Setenv("FAIL_INVALIDINT", "not a int")
	var inv InvalidInt
	Load("FAIL", &inv)
}

func TestInvalidUintShouldPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("this test should panic")
		}
	}()

	type InvalidUint struct {
		InvalidUint uint `env:"invaliduint"`
	}

	os.Setenv("FAIL_INVALIDUINT", "-2")
	var inv InvalidUint
	Load("FAIL", &inv)
}

func TestInvalidBoolShouldPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("this test should panic")
		}
	}()

	type InvalidBool struct {
		InvalidBool bool `env:"invalidbool"`
	}

	os.Setenv("FAIL_INVALIDBOOL", "ture")
	var inv InvalidBool
	Load("FAIL", &inv)
}
