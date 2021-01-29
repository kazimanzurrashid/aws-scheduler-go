package api

import "encoding/json"

type load func(interface{}, interface{}) error

func jsonBasedLoad(s interface{}, d interface{}) error {
	j, err := json.Marshal(s)

	if err != nil {
		return err
	}

	return json.Unmarshal(j, d)
}

var loadStruct load = jsonBasedLoad
