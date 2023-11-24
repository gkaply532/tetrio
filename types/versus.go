package types

import (
	"encoding/json"
	"time"
)

// Handling is the in game handling settings.
type Handling struct {
	ARR      float64 `json:"arr"`
	DAS      float64 `json:"das"`
	DCD      float64 `json:"dcd"`
	SDF      int     `json:"sdf"`
	SafeLock bool    `json:"safelock"`
	Cancel   bool    `json:"cancel"`
}

// VersusStats is the statistics for multiplayer games.
type VersusStats struct {
	APM float64 `json:"apm"` // Attacks per minute
	PPS float64 `json:"pps"` // Pieces per second
	VS  float64 `json:"vs"`  // Versus score
}

type GameRecord struct {
	ID         string          `json:"_id"`
	Stream     string          `json:"stream"`
	ReplayID   string          `json:"replayid"`
	User       PartialUser     `json:"user"`
	TS         time.Time       `json:"ts"`
	IsMulti    bool            `json:"ismulti"`
	EndContext json.RawMessage `json:"endcontext"`
}

type LeagueEndCtx struct {
	ID           string       `json:"id"`
	Username     string       `json:"username"`
	Handling     Handling     `json:"handling"`
	Active       bool         `json:"active"`
	Success      bool         `json:"success"`
	Inputs       int          `json:"inputs"`
	PiecesPlaced int          `json:"piecesplaced"`
	Wins         int          `json:"wins"`
	Points       LeaguePoints `json:"points"`
}

type LeaguePoints struct {
	Primary   int     `json:"primary"`   // Wins
	Secondary float64 `json:"secondary"` // APM
	Tertiary  float64 `json:"tertiary"`  // PPS
	Extra     struct {
		VS float64 `json:"vs"`
	} `json:"extra"`
	SecondaryAvgTracking []float64 `json:"secondaryAvgTracking"` // APM of each round
	TertiaryAvgTracking  []float64 `json:"tertiaryAvgTracking"`  // PPS of each round
	ExtraAvgTracking     struct {
		AggregateStatsVSScore []float64 `json:"aggregatestats___vsscore"`
	} `json:"extraAvgTracking"`
}
