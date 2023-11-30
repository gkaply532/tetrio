package tlhist

import (
	"fmt"
	"strconv"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05.999999999"
const NumFields = 10

type Record struct {
	UserID      string    `json:"user_id"`
	GamesPlayed int       `json:"gamesplayed"`
	GamesWon    int       `json:"gameswon"`
	Rating      float64   `json:"rating"`
	Glicko      float64   `json:"glicko"`
	Rank        string    `json:"rank"`
	APM         float64   `json:"apm"`
	PPS         float64   `json:"pps"`
	VS          float64   `json:"vs"`
	CreatedAt   time.Time `json:"created_at"`
}

func (r *Record) Slice() []string {
	return []string{
		r.UserID,
		strconv.Itoa(r.GamesPlayed),
		strconv.Itoa(r.GamesWon),
		strconv.FormatFloat(r.Rating, 'f', -1, 64),
		strconv.FormatFloat(r.Glicko, 'f', -1, 64),
		r.Rank,
		strconv.FormatFloat(r.APM, 'f', -1, 64),
		strconv.FormatFloat(r.PPS, 'f', -1, 64),
		strconv.FormatFloat(r.VS, 'f', -1, 64),
		r.CreatedAt.Format(TimeLayout),
	}
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func (r *Record) Unslice(fields []string) (err error) {
	s := fields

	if len(s) < NumFields {
		return fmt.Errorf("tlhist: wrong number of fields (expected %d, got %d)", NumFields, len(s))
	}

	defer func() {
		if p := recover(); p != nil {
			e, ok := p.(error)
			if !ok {
				panic(p)
			}

			err = e
		}
	}()

	*r = Record{
		UserID:      s[0],
		GamesPlayed: must(strconv.Atoi(s[1])),
		GamesWon:    must(strconv.Atoi(s[2])),
		Rating:      must(strconv.ParseFloat(s[3], 64)),
		Glicko:      must(strconv.ParseFloat(s[4], 64)),
		Rank:        s[5],
		APM:         must(strconv.ParseFloat(s[6], 64)),
		PPS:         must(strconv.ParseFloat(s[7], 64)),
		VS:          must(strconv.ParseFloat(s[8], 64)),
		CreatedAt:   must(time.Parse(TimeLayout, s[9])),
	}

	return nil
}
