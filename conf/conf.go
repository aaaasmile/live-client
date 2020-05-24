package conf

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ServiceURL            string
	RootURLPattern        string
	DebugVerbose          bool
	IgnoreFatalErrorsInDB bool
	DebugSQL              bool
	DebugParser           bool
	SqliteDBPath          string
	SyncRepo              string
	BeyondComparePath     string
	RemoteServerURL       string
	TempDir               string
	VsCodePath            string
}

var Current = &Config{}

func ReadConfig(configfile string) *Config {
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := toml.DecodeFile(configfile, &Current); err != nil {
		log.Fatal(err)
	}

	return Current
}
