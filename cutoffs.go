package tetrio

import (
	"github.com/gkaply532/tetrio/v2/types"
)

type RankRange struct {
	Rank       string
	Percentile float64
	Top        *types.LeagueUser
	Bottom     *types.LeagueUser
}

// Includes returns wheter TR is in range (Bottom, Top]
func (r RankRange) Includes(TR float64) bool {
	top := r.Top.League.Rating
	bot := r.Bottom.League.Rating
	return top >= TR && TR > bot
}

func Cutoffs(lb []types.LeagueUser) ([]RankRange, error) {
	ranks := []string{"x", "u", "ss", "s+", "s", "s-", "a+", "a", "a-", "b+", "b", "b-", "c+", "c", "c-", "d+", "d"}
	// Percentile values were gathered from here:
	// https://cdn.discordapp.com/attachments/696095920248455229/721851657188671578/league-cheatsheet.png
	percentiles := []float64{0.01, 0.05, 0.11, 0.17, 0.23, 0.3, 0.38, 0.46, 0.54, 0.62, 0.7, 0.78, 0.84, 0.9, 0.95, 0.975, 1}
	cs := make([]RankRange, len(ranks))
	for i := range ranks {
		cs[i] = RankRange{
			Rank:       ranks[i],
			Percentile: percentiles[i],
		}
	}

	// Special case: #1 rated player.
	cs[0].Top = &lb[0]

	for i := 0; i < len(cs)-1; i++ {
		nth := int(float64(len(lb)-1) * cs[i].Percentile)
		cs[i].Bottom = &lb[nth]
		cs[i+1].Top = &lb[nth]
	}

	// Another special case: worst rated player.
	cs[len(cs)-1].Bottom = &lb[len(lb)-1]

	return cs, nil
}
