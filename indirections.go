package tetrio

import "context"

var DefaultClient = NewClient(nil)

func (c *Client) GetUser(ctx context.Context, name string) (*FullUser, error) {
	return c.session.GetUser(ctx, name)
}

func GetUser(name string) (*FullUser, error) {
	return DefaultClient.GetUser(context.Background(), name)
}

func (c *Client) GetCutoffs(ctx context.Context) ([]RankRange, error) {
	return c.session.GetCutoffs(ctx)
}
func GetCutoffs() ([]RankRange, error) {
	return DefaultClient.GetCutoffs(context.Background())
}

func (c *Client) GetLeagueLB(ctx context.Context, country string) ([]LeagueUser, error) {
	return c.session.GetLeagueLB(ctx, country)
}
func GetLeagueLB(country string) ([]LeagueUser, error) {
	return DefaultClient.GetLeagueLB(context.Background(), country)
}

func (c *Client) GetLeagueLBGlobal(ctx context.Context) ([]LeagueUser, error) {
	return c.session.GetLeagueLBGlobal(ctx)
}
func GetLeagueLBGlobal() ([]LeagueUser, error) {
	return DefaultClient.GetLeagueLBGlobal(context.Background())
}

func (c *Client) GetMatches(ctx context.Context, userID string) ([]LeagueRecord, error) {
	return c.session.GetMatches(ctx, userID)
}
func GetMatches(userID string) ([]LeagueRecord, error) {
	return DefaultClient.GetMatches(context.Background(), userID)
}

func (c *Client) SearchUser(ctx context.Context, discordID string) (PartialUser, error) {
	return c.session.SearchUser(ctx, discordID)
}
func SearchUser(discordID string) (PartialUser, error) {
	return DefaultClient.SearchUser(context.Background(), discordID)
}
