package main

import (
	"encoding/json"
	"io/ioutil"
)

func writeJSON(v interface{}, f string) error {
	b, err := json.Marshal(&v)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(f, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func readJSON(f string, v interface{}) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	return nil
}
