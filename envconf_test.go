package envconf

import (
	"os"
	"testing"
)

type Config struct {
	Mongo     MongoConfig `env:"mongo"`
	AppIDList []string    `env:"app_id_list"`
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

	var config Config
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
}
