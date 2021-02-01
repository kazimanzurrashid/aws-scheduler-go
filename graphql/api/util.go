package api

import "encoding/json"

type load func(interface{}, interface{}) error

func jsonBasedLoad(s interface{}, d interface{}) error {
	j, _ := json.Marshal(s)
	return json.Unmarshal(j, d)
}

var loadStruct load = jsonBasedLoad
