package tetrio

import (
	"context"
	"strings"

	"github.com/gkaply532/tetrio/v2/types"
)

func (s Session) GetLeagueLB(ctx context.Context, country string) ([]types.LeagueUser, error) {
	if country == "" {
		return s.GetLeagueLBGlobal(ctx)
	}

	g, err := s.GetLeagueLBGlobal(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]types.LeagueUser, 0)
	for i := range g {
		if strings.EqualFold(g[i].Country, country) {
			result = append(result, g[i])
		}
	}
	return result, nil
}

func (s Session) GetLeagueLBGlobal(ctx context.Context) ([]types.LeagueUser, error) {
	return do[[]types.LeagueUser](ctx, s, "/users/lists/league/all")
}
