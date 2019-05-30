package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type Player struct {
	Id           int64     `json:"id"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Nickname     string    `json:"nickname"`
	Active       bool      `json:"active"`
	RegisteredAt time.Time `json:"registeredAt"`
}

func (p Player) String() string {
	return fmt.Sprintf("Player [%d, %s, %s, %s, %s, %s]", p.Id, p.FirstName, p.LastName, p.Nickname, p.Active, p.RegisteredAt)
}

const (
	getPlayersQuery = "SELECT id, first_name, last_name, nickname, active, registered_at FROM player ORDER BY Id ASC"
	getPlayerQuery  = "SELECT id, first_name, last_name, nickname, active, registered_at FROM player WHERE id = $1"
)

func GetPlayer(dbHandler *sql.DB, id string) []Player {
	players, err := retrievePlayers(dbHandler, getPlayerQuery, id)
	if err != nil {
		logrus.Errorf("Cannot retrieve player (id:%s)", id, err)
		return []Player{}
	}
	return players
}

func GetPlayers(dbHandler *sql.DB) []Player {
	players, err := retrievePlayers(dbHandler, getPlayersQuery)
	if err != nil {
		logrus.Errorf("Cannot retrieve players", err)
		return []Player{}
	}
	return players
}

func retrievePlayers(dbHandler *sql.DB, query string, args ...interface{}) ([]Player, error) {
	var userArray []Player
	err := dbHandler.Ping()
	if err != nil {
		logrus.Errorf("Unable to invoke query: %s!", query, err)
		return nil, err
	}

	var rows *sql.Rows
	if len(args) == 0 {
		rows, err = dbHandler.Query(query)
	} else {
		rows, err = dbHandler.Query(query, args...)
	}
	if err != nil {
		logrus.Errorf("Unable to invoke query: %s!", query, err)
		return nil, err
	}
	defer rows.Close()

	// Convert rows to Player type
	for rows.Next() {
		user, err := convertPlayerRow(rows)
		if err != nil {
			logrus.Warn("Cannot convert player record", err)
		} else {
			userArray = append(userArray, user)
		}
	}
	return userArray, nil
}

func convertPlayerRow(row *sql.Rows) (Player, error) {
	var id int64
	var firstName string
	var lastName string
	var nickname string
	var active bool
	var registeredAt time.Time

	err := row.Scan(&id, &firstName, &lastName, &nickname, &active, &registeredAt)
	switch err {
	case sql.ErrNoRows:
		return Player{}, errors.New("No row to convert!")
	}

	playerStruct := Player{
		Id:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Nickname:     nickname,
		Active:       active,
		RegisteredAt: registeredAt}

	return playerStruct, nil
}
