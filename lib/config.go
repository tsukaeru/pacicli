package lib

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	BaseURL  string
	Username string
	Password string
	Servers  map[string]Server
}

type Server struct {
	Spec          *CreateVe      // xml struct
	Firewall      Firewall      // xml struct
	AutoscaleRule []AutoscaleRule // xml struct
}

func LoadConfig(fpath string, v interface{}) error {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	jsonDecode := func(b []byte, v interface{}) error {
		return json.Unmarshal(b, v)
	}
	tomlDecode := func(b []byte, v interface{}) error {
		_, err := toml.Decode(string(b), v)
		if err != nil {
			return err
		}
		return nil
	}
	switch filepath.Ext(fpath) {
	case "json":
		return jsonDecode(b, v)
	case "toml":
		return tomlDecode(b, v)
	default:
		if len(b) > 0 {
			switch b[0] {
			case '{':
				jsonDecode(b, v)
			default:
				tomlDecode(b, v)
			}
		}
	}
	return nil
}
