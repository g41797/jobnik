package handlers

import (
	"encoding/json"
	"fmt"
)

type serverConfig struct {
	Root string `json:"rootfolder"`
}

func (sc *serverConfig) setDefault() {
	sc.Root = "./"
}

func (sc *serverConfig) unmarshall(js string) error {
	if len(js) == 0 {
		return fmt.Errorf("empty JSON string")
	}

	sc.setDefault()

	return json.Unmarshal([]byte(js), sc)
}
