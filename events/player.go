package events

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	Id int64
	DbId int
	Name string
	SteamID string
	Team string
	Kills int64
	Deaths int64
	Points int
}

type Hit struct {
	Damage int64
	DamageArmor int64
	HitGroup string
	Weapon string
	HitTime string
	Attacker *Player
	Victim *Player
}

type Kill struct {
	KillTime string
	Weapon string
	Attacker *Player
	Victim *Player
}

var Players = make(map[string]*Player)
var Hits []*Hit
var Kills []*Kill

func parsePlayer(name string) *Player {
	var player Player
	rPlayer, _ := regexp.Compile(`^(.*?)<(\d+)><(.*?)><(.*?)>`)
	ParsedPlayer := rPlayer.FindStringSubmatch(name)
	PlayerId, _ := strconv.ParseInt(ParsedPlayer[2], 10, 64)

	if player, ok := Players[ParsedPlayer[1]]; ok {
		player.DbId, player.Points = DbGetPlayer(player)
		player.Team = ParsedPlayer[4]
		return player
	}

	player = Player{Id: PlayerId, Name: ParsedPlayer[1], SteamID: ParsedPlayer[3], Team: ParsedPlayer[4], Kills: 0, Deaths: 0}
	Players[player.Name] = &player

	player.DbId, player.Points = DbGetPlayer(&player)

	return &player
}

func PlayerConnected(message string){
	var player *Player
	if strings.Contains(message, " connected,") {
		rPlayerName, _ := regexp.Compile(`^"(.*?)"`)
		connectedPlayerName := rPlayerName.FindStringSubmatch(message)[1]

		if connectedPlayerName != "Source TV" {
			player = parsePlayer(connectedPlayerName)
			Players[player.Name] = player
			//fmt.Printf("Players: %s\n\n", Players)
			fmt.Printf("[CONNECTED]: %s connected\n", player.Name)
		}
	}
}

func PlayerDisconnected(message string) {
	var player *Player
	if strings.Contains(message, " disconnected (reason ") {
		rPlayerName, _ := regexp.Compile(`^"(.*?)"`)
		disconnectedPlayerName := rPlayerName.FindStringSubmatch(message)[1]

		if disconnectedPlayerName != "Source TV" {
			player = parsePlayer(disconnectedPlayerName)
			delete(Players, player.Name)
			//fmt.Printf("Players: %s\n\n", Players)
			fmt.Printf("[DISCONNECTED]: %s disconnected\n", player.Name)
		}
	}
}

func PlayerTeam(message string) {
	var player *Player
	if strings.Contains(message, " joined team ") {
		rPlayerName, _ := regexp.Compile(`^"(.*?)"`)
		PlayerName := rPlayerName.FindStringSubmatch(message)[1]

		rPlayerTeam, _ := regexp.Compile(`joined team "(.*?)"`)
		PlayerTeam := rPlayerTeam.FindStringSubmatch(message)[1]

		if PlayerName != "Source TV" {
			player = parsePlayer(PlayerName)
			player.Team = PlayerTeam
			Players[player.Name] = player
			//fmt.Printf("Players: %s\n\n", Players)
			fmt.Printf("[TEAM]: %s joined %s\n", player.Name, player.Team)
		}
	}
}

func PlayerKill(message string, pointsKill int, pointsDeath int) {
	var attacker *Player
	var victim *Player
	if strings.Contains(message, " killed ") {
		rKillerName, _ := regexp.Compile(`^"(.*?)"`)
		KillerName := rKillerName.FindStringSubmatch(message)[1]

		rKilledName, _ := regexp.Compile(`killed "(.*?)"`)
		KilledName := rKilledName.FindStringSubmatch(message)[1]
		attacker = parsePlayer(KillerName)
		Players[attacker.Name].Kills++
		fmt.Printf("[KILL]: Going to give %s %d points\n", Players[attacker.Name].Name, pointsKill)
		Players[attacker.Name].Points += pointsKill

		victim = parsePlayer(KilledName)
		Players[victim.Name].Deaths++
		fmt.Printf("[KILL]: Going to give %s %d points\n", Players[victim.Name].Name, pointsDeath)
		Players[victim.Name].Points += pointsDeath

		rWeapon, _ := regexp.Compile(`with "(.*?)"`)
		weapon := rWeapon.FindStringSubmatch(message)[1]

		kill := Kill{KillTime: time.Now().Format("2006-01-02 15:04:05"), Attacker: attacker, Victim: victim, Weapon: weapon}
		Kills = append(Kills, &kill)

		fmt.Printf("[KILL]: Player %s (%d points) killed %s (%d points)\n", attacker.Name, Players[attacker.Name].Points, victim.Name, Players[victim.Name].Points)
		//fmt.Printf("Players: %s\n\n", Players)
	}
}

func PlayerHits(message string, pointsHit int) {
	var attacker *Player
	var victim *Player

	if strings.Contains(message, " attacked ") {
		rAttackerName, _ := regexp.Compile(`^"(.*?)"`)
		AttackerName := rAttackerName.FindStringSubmatch(message)[1]

		rVictimName, _ := regexp.Compile(`attacked "(.*?)"`)
		VictimName := rVictimName.FindStringSubmatch(message)[1]

		rDamage, _ := regexp.Compile(`\(damage "(\d+)"\)`)
		damage, _ := strconv.ParseInt(rDamage.FindStringSubmatch(message)[1], 10, 64)

		rDamageArmor, _ := regexp.Compile(`\(damage_armor "(\d+)"\)`)
		damageArmor, _ := strconv.ParseInt(rDamageArmor.FindStringSubmatch(message)[1], 10, 64)

		rWeapon, _ := regexp.Compile(`with "(.*?)"`)
		weapon := rWeapon.FindStringSubmatch(message)[1]

		rHitGroup, _ := regexp.Compile(`\(hitgroup "(.*?)"\)`)
		hitGroup := rHitGroup.FindStringSubmatch(message)[1]

		attacker = parsePlayer(AttackerName)
		victim = parsePlayer(VictimName)

		Players[attacker.Name].Points += pointsHit

		var hit = Hit{
			HitTime: time.Now().Format("2006-01-02 15:04:05"),
			Attacker: attacker,
			Victim: victim,
			Damage: damage,
			DamageArmor: damageArmor,
			HitGroup: hitGroup,
			Weapon: weapon,
		}

		Hits = append(Hits, &hit)

		fmt.Printf("[HIT]: Gave %s (%d points) %d points\n", attacker.Name, Players[attacker.Name].Points, pointsHit)
		//fmt.Printf("Hits: %s\n\n", Hits)
	}
}
