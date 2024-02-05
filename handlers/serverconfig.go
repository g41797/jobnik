package handlers

import (
	"encoding/json"
	"fmt"
)

type serverConfig struct {
	Addr string `json:"listenonaddr"`
	Root string `json:"rootfolder"`
}

func (sc *serverConfig) setDefault() {
	sc.Addr = ":8080"
	sc.Root = "./"
}

func (sc *serverConfig) unmarshall(js string) error {
	if len(js) == 0 {
		return fmt.Errorf("empty JSON string")
	}

	sc.setDefault()

	return json.Unmarshal([]byte(js), sc)
}
