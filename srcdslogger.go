package main

import (
	"./api"
	"./events"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

type Configuration struct {
	Port int
	ListenIp string
	ApiPort int
	DbHost string
	DbPort int
	DbUser string
	DbPass string
	DbName string
	Points struct {
		Kill int
		Death int
		BombPlanted int
		BombDefused int
		RoundWin int
		RoundLose int
		Hit int
	}
}

var configuration Configuration

type ParsedLogLine struct {
	LogTimestamp time.Time
	SourceServer *net.UDPAddr
	Message string
}

func handleConnection(data string, remote *net.UDPAddr) {
	parsed := parseLogLine(data, remote)
	//fmt.Printf("[%s] [%s] : %s\n", parsed.LogTimestamp, parsed.SourceServer, parsed.Message)
	checkEventsInLog(parsed.Message)
}

func checkEventsInLog(message string) {
	events.PlayerConnected(message)
	events.PlayerDisconnected(message)
	events.PlayerTeam(message)
	events.PlayerKill(message, configuration.Points.Kill, configuration.Points.Death)
	events.PlayerHits(message, configuration.Points.Hit)
	events.ChangeLevel(message)
	events.RoundScore(message)
	events.RoundStart(message)
	events.RoundEnd(message, configuration.Points.RoundWin, configuration.Points.RoundLose)
}

func parseLogLine(line string, remote *net.UDPAddr) ParsedLogLine {
	rLogTimestamp, _ := regexp.Compile(`\d{2}/\d{2}/\d{4} - \d{2}:\d{2}:\d{2}`)
	logTimestamp := rLogTimestamp.FindString(line)
	t, err := time.Parse("01/02/2006 - 15:04:05", logTimestamp)
	if err != nil {
		log.Fatal(err)
	}

	rMessage, _ := regexp.Compile(`: .*`)
	logMessage := strings.TrimPrefix(rMessage.FindString(line), ": ")

	ParsedLogLine := ParsedLogLine{LogTimestamp: t, SourceServer: remote, Message: logMessage}
	return ParsedLogLine
}

func main() {
	file, err := os.Open("config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	laddr := net.UDPAddr{IP: net.ParseIP(configuration.ListenIp), Port: configuration.Port}

	lner, err := net.ListenUDP("udp", &laddr)
	if err != nil {
		log.Fatal(err)
	}
	defer lner.Close()

	events.DbConnection(configuration.DbHost, configuration.DbPort, configuration.DbUser, configuration.DbPass, configuration.DbName)
	api.WebServer(configuration.ListenIp, configuration.ApiPort)
	fmt.Printf("Log Server listening on %s\n", lner.LocalAddr().String())

	for {
		message := make([]byte, 1024)
		rlen, remote, err := lner.ReadFromUDP(message[:])
		if err != nil {
			log.Panic(err)
		}

		data := strings.TrimSpace(string(message[:rlen]))
		handleConnection(data, remote)
	}
}
