package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

const maxMsgLen = 20

func train(words string) {
	w := strings.Split(words, " ")
	modelMap := make(map[string][]string)
	err := readJSON("model.json", &modelMap)
	if err != nil {
		_ = writeJSON(&modelMap, "model.json")
	}
	for i := 0; i < len(w)-1; i++ {
		modelMap[w[i]] = append(modelMap[w[i]], w[i+1])
	}
	modelMap[w[len(w)-1]] = []string{}
	err = writeJSON(&modelMap, "model.json")
	if err != nil {
		log.Fatal("error writing model.json: ", err)
	}
}

func generate() string {
	rand.Seed(time.Now().UnixNano())
	var currentWord string
	modelMap := make(map[string][]string)
	err := readJSON("model.json", &modelMap)
	if err != nil {
		log.Fatal("error reading model.json: ", err)
	}
	for k := range modelMap {
		currentWord = k
		break
	}
	result := currentWord
	for i := 0; i < maxMsgLen; i++ {
		possibilities := modelMap[currentWord]
		pLen := len(possibilities)
		if pLen <= 0 {
			continue
		}
		next := possibilities[rand.Intn(pLen)]
		currentWord = next
		result += fmt.Sprintf(" %s ", next)
	}
	return result
}
