package types

import "time"

type PartialUser struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
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

type FullTLStats struct {
	PartialTLStats

	NextRank       string `json:"next_rank"`
	PrevRank       string `json:"prev_rank"`
	PercentileRank string `json:"percentile_rank"`

	// The #1 ranked player has Standing 1. (not 0)
	Standing      int     `json:"standing"`
	StandingLocal int     `json:"standing_local"`
	NextAt        int     `json:"next_at"`
	PrevAt        int     `json:"prev_at"`
	Percentile    float64 `json:"percentile"`
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
