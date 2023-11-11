package tetrio

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
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

var ErrNonMultiRecord = errors.New("tetrio: not a multiplayer record")
var ErrAmbiguousRecord = errors.New("tetrio: record has no winner or loser")

type LeagueRecord struct {
	ReplayID  string
	TS        time.Time
	IsForfeit bool
	Winner    LeaguePlayer
	Loser     LeaguePlayer
}

type LeaguePlayer struct {
	User         PartialUser
	Wins         int
	Inputs       int
	PiecesPlaced int
	Handling     Handling
	Stats        VersusStats
	RoundStats   []VersusStats
}

func (s Session) GetMatches(ctx context.Context, userID string) ([]LeagueRecord, error) {
	return send[[]LeagueRecord](ctx, s, "/streams/league_userrecent_"+url.PathEscape(userID))
}

func (g *LeagueRecord) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var rec gameRecord
	err := json.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	if !rec.IsMulti {
		return ErrNonMultiRecord
	}
	var cs [2]leagueEndCtx
	err = json.Unmarshal(rec.EndContext, &cs)
	if err != nil {
		return err
	}

	if cs[0].Success == cs[1].Success {
		return ErrAmbiguousRecord
	}

	g.ReplayID = rec.ReplayID
	g.TS = rec.TS
	g.IsForfeit = cs[0].Active != cs[1].Active

	for i := 0; i < 2; i++ {
		apms := cs[i].Points.SecondaryAvgTracking
		ppss := cs[i].Points.TertiaryAvgTracking
		vss := cs[i].Points.ExtraAvgTracking.AggregateStatsVSScore

		roundCount := max(len(apms), len(ppss), len(vss))
		rs := make([]VersusStats, roundCount)
		for j, apm := range apms {
			rs[j].APM = apm
		}
		for j, pps := range ppss {
			rs[j].PPS = pps
		}
		for j, vs := range vss {
			rs[j].VS = vs
		}

		p := LeaguePlayer{
			User: PartialUser{
				ID:       cs[i].ID,
				Username: cs[i].Username,
			},
			Handling:     cs[i].Handling,
			Inputs:       cs[i].Inputs,
			PiecesPlaced: cs[i].PiecesPlaced,
			Wins:         cs[i].Wins,
			Stats: VersusStats{
				APM: cs[i].Points.Secondary,
				PPS: cs[i].Points.Tertiary,
				VS:  cs[i].Points.Extra.VS,
			},
			RoundStats: rs,
		}
		if cs[i].Success {
			g.Winner = p
		} else {
			g.Loser = p
		}
	}
	return nil
}

type gameRecord struct {
	ID         string          `json:"_id"`
	Stream     string          `json:"stream"`
	ReplayID   string          `json:"replayid"`
	User       PartialUser     `json:"user"`
	TS         time.Time       `json:"ts"`
	IsMulti    bool            `json:"ismulti"`
	EndContext json.RawMessage `json:"endcontext"`
}

type leagueEndCtx struct {
	ID           string       `json:"id"`
	Username     string       `json:"username"`
	Handling     Handling     `json:"handling"`
	Active       bool         `json:"active"`
	Success      bool         `json:"success"`
	Inputs       int          `json:"inputs"`
	PiecesPlaced int          `json:"piecesplaced"`
	Wins         int          `json:"wins"`
	Points       leaguePoints `json:"points"`
}

type leaguePoints struct {
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
