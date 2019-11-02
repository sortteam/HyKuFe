package horovodjob

import (
	"encoding/json"
)

func DefinitionToJson(instance interface{}) []byte{
	jsonByte, err := json.Marshal(instance)
	if err != nil {

	}
	return jsonByte
}