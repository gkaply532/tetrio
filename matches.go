package tetrio

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"time"

	"github.com/gkaply532/tetrio/v2/types"
)

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
	User         types.PartialUser
	Wins         int
	Inputs       int
	PiecesPlaced int
	Handling     types.Handling
	Stats        types.VersusStats
	RoundStats   []types.VersusStats
}

func (s Session) GetMatches(ctx context.Context, userID string) ([]LeagueRecord, error) {
	return do[[]LeagueRecord](ctx, s, "/streams/league_userrecent_"+url.PathEscape(userID))
}

func (g *LeagueRecord) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var rec types.GameRecord
	err := json.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	if !rec.IsMulti {
		return ErrNonMultiRecord
	}
	var cs [2]types.LeagueEndCtx
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
		rs := make([]types.VersusStats, roundCount)
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
			User: types.PartialUser{
				ID:       cs[i].ID,
				Username: cs[i].Username,
			},
			Handling:     cs[i].Handling,
			Inputs:       cs[i].Inputs,
			PiecesPlaced: cs[i].PiecesPlaced,
			Wins:         cs[i].Wins,
			Stats: types.VersusStats{
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
