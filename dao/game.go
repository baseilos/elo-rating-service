package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type Game struct {
	Id          int64     `json:"id"`
	WhitePlayer int64     `json:"whitePlayer"`
	BlackPlayer int64     `json:"blackPlayer"`
	Result      int       `json:"result"`
	PlayedAt    time.Time `json:"playedAt"`
}

func (p Game) String() string {
	return fmt.Sprintf("Game [%d, %d, %d, %d, %s]", p.Id, p.WhitePlayer, p.BlackPlayer, p.Result, p.PlayedAt)
}

const (
	getGamesQuery = "SELECT id, player_white, player_black, played_at, result FROM game ORDER BY played_at DESC"
)

func GetGames(dbHandler *sql.DB) []Game {
	return retrieveGames(dbHandler, getGamesQuery)
}

func retrieveGames(dbHandler *sql.DB, query string, args ...interface{}) []Game {
	var gameArray []Game
	err := dbHandler.Ping()
	if err != nil {
		logrus.Errorf("Unable to invoke query: %s!", query, err)
	}

	var rows *sql.Rows
	if len(args) == 0 {
		rows, err = dbHandler.Query(query)
	} else {
		rows, err = dbHandler.Query(query, args...)
	}
	if err != nil {
		logrus.Errorf("Unable to invoke query: %s!", query, err)
		return gameArray
	}
	defer rows.Close()

	// Convert rows to Player type
	for rows.Next() {
		user, err := convertGameRow(rows)
		if err != nil {
			logrus.Warn("Cannot convert player record", err)
		} else {
			gameArray = append(gameArray, user)
		}
	}
	return gameArray
}

func convertGameRow(row *sql.Rows) (Game, error) {
	var id int64
	var whitePlayer int64
	var blackPlayer int64
	var result int
	var playedAt time.Time

	err := row.Scan(&id, &whitePlayer, &blackPlayer, &playedAt, &result)
	switch err {
	case sql.ErrNoRows:
		return Game{}, errors.New("No row to convert!")
	}

	gameStruct := Game{
		Id:          id,
		WhitePlayer: whitePlayer,
		BlackPlayer: blackPlayer,
		Result:      result,
		PlayedAt:    playedAt}

	return gameStruct, nil
}
