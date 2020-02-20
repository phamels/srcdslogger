package events

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

var db *sql.DB
var err error

func DbConnection(DbHost string, DbPort int, DbUser string, DbPass string, DbName string) (*sql.DB, error) {
	db, err = sql.Open("mysql", DbUser+":"+DbPass+"@tcp("+DbHost+":"+strconv.Itoa(DbPort)+")/"+DbName+"?parseTime=true")

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Printf("Connected to db server %s:%d\n", DbHost, DbPort)

	return db, nil
}

func DbGetPlayer(player *Player) (int, int) {
	id := 0
	points := int(0)

	if player.SteamID == "BOT" {
		err := db.QueryRow(`SELECT id, points FROM players WHERE name = ? AND steamid = ?`, player.Name, player.SteamID).Scan(&id, &points)
		if err != nil {
			if err == sql.ErrNoRows {
				sqlStmt := `insert into players (name, steamid, points) values (?, ?, ?)`
				res, err := db.Exec(sqlStmt, player.Name, player.SteamID, points)
				if err != nil {
					log.Println(err)
				}
				lid, err := res.LastInsertId()
				if err != nil {
					log.Println(err)
				}
				id = int(lid)
			} else {
				log.Println(err)
			}
		}
	} else {
		err := db.QueryRow(`SELECT id, points FROM players WHERE steamid = ?`, player.SteamID).Scan(&id, &points)
		if err != nil {
			if err == sql.ErrNoRows {
				sqlStmt := `insert into players (name, steamid, points) values (?, ?, ?)`
				res, err := db.Exec(sqlStmt, player.Name, player.SteamID, points)
				if err != nil {
					log.Println(err)
				}
				lid, err := res.LastInsertId()
				if err != nil {
					log.Println(err)
				}
				id = int(lid)
			} else {
				log.Println(err)
			}
		}
	}

	return id, points
}

func DbGetMapId(mapname string) int {
	id := 0

	err := db.QueryRow(`SELECT id FROM maps WHERE map = ?`, mapname).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			sqlStmt := `insert into maps (map) values (?)`
			res, err := db.Exec(sqlStmt, mapname)
			if err != nil {
				log.Println(err)
			}
			lid, err := res.LastInsertId()
			if err != nil {
				log.Println(err)
			}
			id = int(lid)
		} else {
			log.Println(err)
		}
	}

	return id
}

func DbGetHitGroupId(hitgroup string) int {
	id := 0

	err := db.QueryRow(`SELECT id FROM hitgroups WHERE hitgroup = ?`, hitgroup).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			sqlStmt := `insert into hitgroups (hitgroup) values (?)`
			res, err := db.Exec(sqlStmt, hitgroup)
			if err != nil {
				log.Println(err)
			}
			lid, err := res.LastInsertId()
			if err != nil {
				log.Println(err)
			}
			id = int(lid)
		} else {
			log.Println(err)
		}
	}

	return id
}

func DbGetWeaponId(weapon string) int {
	id := 0

	err := db.QueryRow(`SELECT id FROM weapons WHERE weapon = ?`, weapon).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			sqlStmt := `insert into weapons (weapon) values (?)`
			res, err := db.Exec(sqlStmt, weapon)
			if err != nil {
				log.Println(err)
			}
			lid, err := res.LastInsertId()
			if err != nil {
				log.Println(err)
			}
			id = int(lid)
		} else {
			log.Println(err)
		}
	}

	return id
}

func DbInsertRoundResult(round *Round) {
	roundId := 0
	mapId := DbGetMapId(round.Level)

	sqlStmt := `insert into rounds (ct_score, t_score, outcome, winner, map_id, start_time, end_time) values (?, ?, ?, ?, ?, ?, ?)`
	res, err := db.Exec(sqlStmt, round.CTScore, round.TScore, round.Outcome, round.Winner, mapId, round.StartTime, round.EndTime)
	if err != nil {
		log.Println(err)
	}
	lid, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	roundId = int(lid)

	if roundId != 0 {
		for _, player := range round.Players {
			sqlStmt := `insert into player_round (player_id, round_id) values (?, ?)`
			_, err := db.Exec(sqlStmt, player.DbId, roundId)
			if err != nil {
				log.Println(err)
			}
		}

		for _, kill := range round.Kills {
			weaponId := DbGetWeaponId(kill.Weapon)
			sqlStmt := `insert into kills (round_id, attacker_id, victim_id, weapon_id, kill_time) values (?, ?, ?, ?, ?)`
			_, err := db.Exec(sqlStmt, roundId, kill.Attacker.DbId, kill.Victim.DbId, weaponId, kill.KillTime)
			if err != nil {
				log.Println(err)
			}
		}

		for _, hit := range round.Hits {
			hitgroupId := DbGetHitGroupId(hit.HitGroup)
			weaponId := DbGetWeaponId(hit.Weapon)
			sqlStmt := `insert into hits (round_id, hitgroup_id, weapon_id, damage, damage_armor, attacker_id, victim_id, hit_time) values (?, ?, ?, ?, ?, ?, ?, ?)`
			_, err := db.Exec(sqlStmt, roundId, hitgroupId, weaponId, hit.Damage, hit.DamageArmor, hit.Attacker.DbId, hit.Victim.DbId, hit.HitTime)
			if err != nil {
				log.Println(err)
			}
		}

		for _, player := range round.Players {
			sqlStmt := `update players set points = ? where id = ?`
			fmt.Printf("[ROUND END]: Setting points %d for player %s (%s)\n", player.Points, player.Name, player.Team)
			_, err := db.Exec(sqlStmt, player.Points, player.DbId)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
