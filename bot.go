package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"
	"regexp"
	"strings"
)

const (
	twitchIRC = "irc.chat.twitch.tv:6667"
)

var (
	messageRegex = *regexp.MustCompile(`^* PRIVMSG.* :(.*)$`)
	username     = os.Getenv("USERNAME")
	channel      = os.Getenv("STRIMER")
	oauth        = os.Getenv("OAUTH")
	verbose      bool
)

type twitchBot struct {
	username string
	channel  string
	oauth    string
	conn     net.Conn
	messages chan string
}

func initBot(u, ch, o string) *twitchBot {
	c, err := net.Dial("tcp", twitchIRC)
	if err != nil {
		log.Fatal("can't connect to IRC server: ", err)
	}
	bot := &twitchBot{
		username: u,
		channel:  ch,
		oauth:    o,
		conn:     c,
		messages: make(chan string),
	}
	go bot.receive()
	bot.send("PASS " + o)
	bot.send("NICK " + u)
	bot.join()
	return bot
}

func (bot *twitchBot) receive() {
	tp := textproto.NewReader(bufio.NewReader(bot.conn))
	for {
		message, err := tp.ReadLine()
		if verbose {
			fmt.Println("> ", message)
		}
		if err != nil {
			log.Fatal("lost connection to server: ", err)
		}
		bot.messages <- message
	}
}

func (bot *twitchBot) send(m string) {
	_, err := bot.conn.Write([]byte(m + "\r\n"))
	if verbose {
		fmt.Println("< ", m)
	}
	if err != nil {
		log.Fatal("lost connection to server: ", err)
	}
}

func (bot *twitchBot) join() {
	message := fmt.Sprintf("JOIN #%s", bot.channel)
	bot.send(message)
}

func (bot *twitchBot) chat(m string) {
	message := fmt.Sprintf("PRIVMSG #%s :%s", bot.channel, m)
	bot.send(message)
}

func main() {
	v := flag.Bool("v", false, "verbose chat (print chat to stdout)")
	l := flag.Bool("l", false, "learning mode (only listen to chat)")
	flag.Parse()
	verbose = *v
	var (
		trainData string
		msgCount  int
	)
	bot := initBot(username, channel, oauth)
	fmt.Println("Connected to", twitchIRC)
	for m := range bot.messages {
		if strings.HasPrefix(m, "PING") {
			bot.send("PONG:tmi.twitch.tv")
			continue
		}
		matches := messageRegex.FindStringSubmatch(m)
		if len(matches) <= 1 {
			continue
		}
		trainData += fmt.Sprintf("%s ", matches[1])
		msgCount++
		if msgCount >= 150 {
			train(trainData)
			if !*l {
				chatMsg := generate()
				bot.chat(chatMsg)
			}
			trainData = ""
			msgCount = 0
		}
	}
}
