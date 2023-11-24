package tetrio

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gkaply532/tetrio/v2/types"
)

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

type SearchError uint64

func (e SearchError) Error() string {
	return fmt.Sprintf("tetrio: no linked user found for discord snowflake: %d", uint64(e))
}

func (s Session) GetUser(ctx context.Context, name string) (*types.FullUser, error) {
	if len(name) > 60 {
		return nil, errors.New("tetrio: username too long")
	}

	endpoint := "/users/" + url.PathEscape(strings.ToLower(name))
	user, err := do[*types.FullUser](ctx, s, endpoint)
	if err != nil {
		var tetrioErr Error
		if errors.As(err, &tetrioErr) && strings.Contains(
			strings.ToLower(string(tetrioErr)),
			"no such user",
		) {
			return nil, NoUserError{
				UID:     name,
				wrapped: err,
			}
		}
		return nil, err
	}

	return user, nil
}

func (s Session) SearchUser(ctx context.Context, discordID uint64) (types.PartialUser, error) {
	user, err := do[types.PartialUser](ctx, s, "/users/search/"+url.PathEscape(strconv.FormatUint(discordID, 10)))
	if err != nil {
		return types.PartialUser{}, err
	}
	if user.ID == "" {
		return types.PartialUser{}, SearchError(discordID)
	}
	return user, nil
}
