package types

type LeagueUser struct {
	PartialUser

	Role      string         `json:"role"`
	XP        float64        `json:"xp"`
	Country   string         `json:"country"`
	Supporter bool           `json:"supporter"`
	Verified  bool           `json:"verified"`
	League    PartialTLStats `json:"league"`
}

type PartialTLStats struct {
	GamesPlayed int     `json:"gamesplayed"`
	GamesWon    int     `json:"gameswon"`
	Rating      float64 `json:"rating"`
	Rank        string  `json:"rank"`
	BestRank    string  `json:"bestrank"`
	Glicko      float64 `json:"glicko"`
	RD          float64 `json:"rd"`
	APM         float64 `json:"apm"`
	PPS         float64 `json:"pps"`
	VS          float64 `json:"vs"`
	Decaying    bool    `json:"decaying"`
}
