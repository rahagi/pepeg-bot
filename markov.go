package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const maxMsgLen = 20

var (
	fullText    string
	lookupTable = make(map[string][]string)
)

func train(words string) {
	fullText += words
	w := strings.Split(fullText, " ")
	for i := 0; i < len(w)-1; i++ {
		lookupTable[w[i]] = append(lookupTable[w[i]], w[i+1])
	}
	lookupTable[w[len(w)-1]] = []string{}
}

func generate() string {
	rand.Seed(time.Now().UnixNano())
	words := strings.Split(fullText, " ")
	currentWord := words[rand.Intn(len(words))]
	result := currentWord
	for i := 0; i < maxMsgLen; i++ {
		possibilities := lookupTable[currentWord]
		pLen := len(possibilities)
		if pLen <= 0 {
			continue
		}
		next := possibilities[rand.Intn(pLen)]
		currentWord = next
		result += fmt.Sprintf(" %s ", next)
	}
	fullText = ""
	return result
}
