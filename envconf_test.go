package envconf

import (
	"os"
	"testing"
)

type Config struct {
	Mongo     MongoConfig `env:"mongo"`
	AppIDList []string    `env:"app_id_list"`
	Replicas  uint        `env:"replicas"`

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

func TestLoad(t *testing.T) {
	os.Setenv("TEST_MONGO_NODES", "www.example.com")
	os.Setenv("TEST_MONGO_PORT", "332")
	os.Setenv("TEST_MONGO_DEBUG", "false")
	os.Setenv("TEST_APP_ID_LIST", " aa, bb ,cc ,dd")
	os.Setenv("TEST_REPLICAS", "3")

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

	if config.Mongo.Port != 332 {
		t.FailNow()
	}

	if config.Mongo.Nodes != "www.example.com" {
		t.FailNow()
	}

	if config.Mongo.Database != "" {
		t.FailNow()
	}

	if config.Mongo.Debug {
		t.FailNow()
	}

	expAppIDList := []string{"aa", "bb", "cc", "dd"}

	if len(config.AppIDList) != len(expAppIDList) {
		t.FailNow()
	}

	for i := range config.AppIDList {
		if config.AppIDList[i] != expAppIDList[i] {
			t.FailNow()
		}
	}

	if config.Replicas != uint(3) {
		t.FailNow()
	}

	if config.unexported != initConfig.unexported {
		t.FailNow()
	}

	if config.UnTagged != initConfig.UnTagged {
		t.FailNow()
	}

	if config.UnLoadedInt != initConfig.UnLoadedInt {
		t.FailNow()
	}

	if config.UnLoadedUint != initConfig.UnLoadedUint {
		t.FailNow()
	}

	if config.UnLoadedBool != initConfig.UnLoadedBool {
		t.FailNow()
	}

	if config.UnLoadedStr != initConfig.UnLoadedStr {
		t.FailNow()
	}

	if len(config.UnLoadedStrSlice) != len(initConfig.UnLoadedStrSlice) {
		t.FailNow()
	}

	for i := range config.UnLoadedStrSlice {
		if config.UnLoadedStrSlice[i] != initConfig.UnLoadedStrSlice[i] {
			t.FailNow()
		}
	}
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
