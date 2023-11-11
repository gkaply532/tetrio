package tetrio

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var leagueDumpPath = wellKnownTempPath("tetrio-leaderboard-dump.json")
var leagueDumpLock sync.Mutex

type PartialTLStats struct {
	// Standing isn't returned by the API for the TETRA LEAGUE leaderboard list.
	// It is filled in by the library.
	//
	// The #1 ranked player has Standing 1. (not 0)
	Standing int `json:"standing"`

	GamesPlayed int      `json:"gamesplayed"`
	GamesWon    int      `json:"gameswon"`
	Rating      float64  `json:"rating"`
	Rank        string   `json:"rank"`
	BestRank    *string  `json:"bestrank,omitempty"`
	Glicko      float64  `json:"glicko"`
	RD          float64  `json:"rd"`
	APM         float64  `json:"apm"`
	PPS         float64  `json:"pps"`
	VS          *float64 `json:"vs,omitempty"`
	Decaying    *bool    `json:"decaying,omitempty"`
}

type LeagueUser struct {
	PartialUser

	Role      string         `json:"role"`
	XP        *float64       `json:"xp"`
	Country   *string        `json:"country"`
	Supporter bool           `json:"supporter"`
	Verified  *bool          `json:"verified"`
	League    PartialTLStats `json:"league"`
}

type TenchiDump struct {
	Success bool         `json:"success"`
	Users   []LeagueUser `json:"users"`
	TS      time.Time    `json:"ts"`
}

type UserSnapshot struct {
	TS   time.Time  `json:"ts"`
	User LeagueUser `json:"user"`
}

type RankRange struct {
	Rank       string
	Percentile float64
	Top        *LeagueUser
	Bottom     *LeagueUser
}

// Includes returns wheter TR is in range (Bottom, Top]
func (r RankRange) Includes(TR float64) bool {
	top := r.Top.League.Rating
	bot := r.Bottom.League.Rating
	return top >= TR && TR > bot
}

func wellKnownTempPath(name string) string {
	return filepath.Join(os.TempDir(), name)
}

// TODO(gkm): The results of these exported functions should be cached because
// they are expensive to run.

func (s Session) GetLeagueLB(ctx context.Context, country string) ([]LeagueUser, error) {
	g, err := s.GetLeagueLBGlobal(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]LeagueUser, 0)
	for i := range g {
		var c string
		if g[i].Country != nil {
			c = *g[i].Country
		}
		if strings.EqualFold(c, country) {
			result = append(result, g[i])
		}
	}
	return result, nil
}

func (s Session) GetLeagueLBGlobal(ctx context.Context) ([]LeagueUser, error) {
	leagueDumpLock.Lock()
	defer leagueDumpLock.Unlock()
	file, err := os.OpenFile(leagueDumpPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	isEmpty := info.Size() == 0
	expired := isEmpty

	if !isEmpty {
		resp, err := parseResponse(file)
		if err != nil {
			return nil, fmt.Errorf("tetrio: can't parse the league dump file: %w", err)
		}
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		expired = resp.Cache.Until.Before(time.Now())
	}

	jsonReader := io.Reader(file)

	if expired {
		err = file.Truncate(0)
		if err != nil {
			return nil, err
		}

		resp, err := s.makeRequest(ctx, "/users/lists/league/all")
		if err != nil {
			return nil, err
		}
		defer resp.Close()

		jsonReader = io.TeeReader(resp, file)
	}

	users, err := parseDataFromJSON[[]LeagueUser](jsonReader)
	if err != nil {
		return nil, err
	}
	for i := range users {
		users[i].League.Standing = i + 1
	}

	return users, nil
}

func cutoffsSlice() []RankRange {
	ranks := []string{"x", "u", "ss", "s+", "s", "s-", "a+", "a", "a-", "b+", "b", "b-", "c+", "c", "c-", "d+", "d"}
	// Percentile values were gathered from here:
	// https://cdn.discordapp.com/attachments/696095920248455229/721851657188671578/league-cheatsheet.png
	percentiles := []float64{0.01, 0.05, 0.11, 0.17, 0.23, 0.3, 0.38, 0.46, 0.54, 0.62, 0.7, 0.78, 0.84, 0.9, 0.95, 0.975, 1}
	result := make([]RankRange, len(ranks))
	for i := range ranks {
		result[i] = RankRange{
			Rank:       ranks[i],
			Percentile: percentiles[i],
		}
	}
	return result
}

func (s Session) GetCutoffs(ctx context.Context) ([]RankRange, error) {
	lb, err := s.GetLeagueLBGlobal(ctx)
	if err != nil {
		return nil, err
	}

	cutoffs := cutoffsSlice()

	// Special case: #1 rated player.
	cutoffs[0].Top = &lb[0]

	for i := 0; i < len(cutoffs)-1; i++ {
		nth := int(float64(len(lb)-1) * cutoffs[i].Percentile)
		cutoffs[i].Bottom = &lb[nth]
		cutoffs[i+1].Top = &lb[nth]
	}

	// Another special case: worst rated player.
	cutoffs[len(cutoffs)-1].Bottom = &lb[len(lb)-1]

	return cutoffs, nil
}
