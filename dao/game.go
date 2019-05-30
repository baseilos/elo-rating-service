package dao

import (
	"database/sql"
	"encoding/json"
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
	getGamesQuery  = "SELECT id, player_white, player_black, played_at, result FROM game ORDER BY played_at DESC"
	storeGameQuery = "INSERT INTO game (player_white, player_black, played_at, result) VALUES($1,$2,$3, $4)"
)

func FromJson(jsonString string) (*Game, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		logrus.Errorf("Cannot unmarshal game information: %s", jsonString, err)
		return nil, err
	}

}

func GetGames(dbHandler *sql.DB) []Game {
	games, err := retrieveGames(dbHandler, getGamesQuery)
	if err != nil {
		logrus.Errorf("Cannot retrieve games!", err)
		return []Game{}
	}
	return games
}

func AddGame(dbHandler *sql.DB, game Game) Game {
	// TOOD: Finish
	return game
}

func retrieveGames(dbHandler *sql.DB, query string, args ...interface{}) ([]Game, error) {
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
	var gameArray = []Game{}
	for rows.Next() {
		game, err := convertGameRow(rows)
		if err != nil {
			logrus.Warn("Cannot convert player record", err)
		} else {
			gameArray = append(gameArray, game)
		}
	}
	return gameArray, nil
}

func storeGame(dbHandler *sql.DB, game Game) (int64, error) {
	err := dbHandler.Ping()
	if err != nil {
		logrus.Errorf("Cannot store game: %s. No database connection!", game, err)
		return -1, err
	}

	// Open new transaction and prepare SQL statement
	tx, err := dbHandler.Begin()
	if err != nil {
		logrus.Errorf("Cannot store game: %s. Cannot obtain transaction!", game, err)
		return -1, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(storeGameQuery)
	if err != nil {
		logrus.Errorf("Cannot store game: %s. Cannot obtain prepared statement!", game, err)
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(game.WhitePlayer, game.BlackPlayer, game.PlayedAt, game.Result)
	if err != nil {
		logrus.Errorf("Cannot store game: %s. Insert failed!", game, err)
		return -1, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		logrus.Errorf("Cannot store game: %s. Cannot obtain last inserted id!", game, err)
		return -1, err
	}
	return lastId, nil

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

func convertGameRowFromMap(data map[string]interface{}) (Game, error) {
	gameStruct := Game{
		WhitePlayer: data["whitePlayer"].(int64),
		BlackPlayer: data["blackPlayer"].(int64),
		Result:      data["result"].(int),
		PlayedAt:    data["playedAt"].(time.Time)}

	return gameStruct, nil
}
