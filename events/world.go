package events

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Round struct {
	StartTime string
	EndTime string
	CTScore int64
	TScore int64
	Winner string
	Outcome string
	Level string
	Players map[string]*Player
	Hits []*Hit
	Kills []*Kill
}

type RoundOutcome struct {
	Winner string
	Outcome string
	CTScore int64
	TScore int64
}

var Rounds []*Round
var roundOutcome RoundOutcome
var Level = ""
var round Round

func ChangeLevel(message string) {
	if strings.Contains(message, "Started map") {
		rLevel, _ := regexp.Compile(`Started map "(.*?)"`)
		Level = rLevel.FindStringSubmatch(message)[1]
		Rounds = Rounds[:0]
	}
}

func RoundScore(message string) {
	if strings.Contains(message, "\"Terrorists_Win\"") || strings.Contains(message, "\"Target_Bombed\"") {
		rOutcome, _ := regexp.Compile(`triggered "(.*?)"`)
		outcome := rOutcome.FindStringSubmatch(message)[1]

		rCTScore, _ := regexp.Compile(`\(CT "(\d+)"\)`)
		ctscore, _ := strconv.ParseInt(rCTScore.FindStringSubmatch(message)[1], 10, 64)

		rTScore, _ := regexp.Compile(`\(T "(\d+)"\)`)
		tscore, _ := strconv.ParseInt(rTScore.FindStringSubmatch(message)[1], 10, 64)

		roundOutcome = RoundOutcome{Winner: "TERRORIST", Outcome: outcome, CTScore: ctscore, TScore: tscore}
	}

	if strings.Contains(message, "\"CTs_Win\"") || strings.Contains(message, "\"Bomb_Defused\"") {
		rOutcome, _ := regexp.Compile(`triggered "(.*?)"`)
		outcome := rOutcome.FindStringSubmatch(message)[1]

		rCTScore, _ := regexp.Compile(`\(CT "(\d+)"\)`)
		ctscore, _ := strconv.ParseInt(rCTScore.FindStringSubmatch(message)[1], 10, 64)

		rTScore, _ := regexp.Compile(`\(T "(\d+)"\)`)
		tscore, _ := strconv.ParseInt(rTScore.FindStringSubmatch(message)[1], 10, 64)

		roundOutcome = RoundOutcome{Winner: "CT", Outcome: outcome, CTScore: ctscore, TScore: tscore}
	}
}

func RoundStart(message string) {
	if strings.Contains(message, "World triggered \"Round_Start\"") {
		roundOutcome = RoundOutcome{}
		fmt.Printf(">>> New round\n\n")
		round.StartTime = time.Now().Format("2006-01-02 15:04:05")
	}
}

func RoundEnd(message string, pointsWin int, pointsLose int) {
	if strings.Contains(message, "World triggered \"Round_End\"") {
		round.EndTime = time.Now().Format("2006-01-02 15:04:05")
		round.CTScore = roundOutcome.CTScore
		round.TScore = roundOutcome.TScore
		round.Outcome = roundOutcome.Outcome
		round.Winner = roundOutcome.Winner
		round.Level = Level
		round.Players = Players
		round.Hits = Hits
		round.Kills = Kills

		for _, player := range round.Players {
			if player.Team == round.Winner {
				player.Points += pointsWin
			} else if player.Team == "Spectator" {

			} else {
				player.Points += pointsLose
			}
		}

		Rounds = append(Rounds, &round)

		if round.StartTime != "" {
			DbInsertRoundResult(&round)
		}

		fmt.Printf(">>> Round ended\n\n")
		fmt.Printf("%+v\n\n", round)

		round.StartTime = ""
		round.EndTime = ""
		round.Outcome = ""
		Hits = Hits[:0]
		Kills = Kills[:0]
	}
}
