package tetrio

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type PartialUser struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
}

type Badge struct {
	ID    string    `json:"id"`
	Label string    `json:"label"`
	TS    time.Time `json:"ts"`
}

type DiscordInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type FullTLStats struct {
	PartialTLStats

	NextRank       string `json:"next_rank"`
	PrevRank       string `json:"prev_rank"`
	PercentileRank string `json:"percentile_rank"`

	StandingLocal int     `json:"standing_local"`
	NextAt        int     `json:"next_at"`
	PrevAt        int     `json:"prev_at"`
	Percentile    float64 `json:"percentile"`
}

type FullUser struct {
	LeagueUser
	League FullTLStats `json:"league"`

	TS        time.Time `json:"ts"`
	BotMaster string    `json:"botmaster"`
	Badges    []Badge   `json:"badges"`

	GamesPlayed int     `json:"gamesplayed"`
	GamesWon    int     `json:"gameswon"`
	GameTime    float64 `json:"gametime"`

	BadStanding   bool `json:"badstanding"`
	SupporterTier int  `json:"supporter_tier"`

	AvatarRevision int64 `json:"avatar_revision"`
	BannerRevision int64 `json:"banner_revision"`

	Bio         string `json:"bio"`
	FriendCount int    `json:"friend_count"`
	Connections struct {
		Discord DiscordInfo
	} `json:"connections"`
	Distinguishment struct {
		Type string
	} `json:"distinguishment"`
}

type NoUserError struct {
	UID     string
	wrapped error
}

func (e NoUserError) Error() string {
	return fmt.Sprintf("tetrio: no such user: %q", e.UID)
}

func (e NoUserError) Unwrap() error {
	return e.wrapped
}

func GetUser(name string) (*FullUser, error) {
	if len(name) > 60 {
		return nil, errors.New("tetrio: long username")
	}
	endpoint := "/users/" + url.PathEscape(strings.ToLower(name))
	user, err := send[*FullUser](endpoint)
	var tetrioErr Error
	if errors.As(err, &tetrioErr) && strings.Contains(
		strings.ToLower(string(tetrioErr)),
		"no such user",
	) {
		err = NoUserError{
			UID:     name,
			wrapped: err,
		}
	}
	return user, err
}

// SearchUser returns the PartialUser that has their Discord account linked with
// TETR.IO. returns a zero PartialUser if no user is found.
func SearchUser(discordID string) (PartialUser, error) {
	return send[PartialUser]("/users/search/" + url.PathEscape(discordID))
}
