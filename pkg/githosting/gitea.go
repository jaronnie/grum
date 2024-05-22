package githosting

import (
	"context"
	"github.com/jaronnie/genius"
	"github.com/spf13/cast"

	"github.com/jaronnie/restc"
)

type Gitea struct {
	Config Config
	Client restc.Interface
}

func (g *Gitea) GetUserInfo() (*UserInfo, error) {
	response, err := g.Client.Get().SubPath("/api/v1/user").Do(context.Background()).RawResponse()
	if err != nil {
		return nil, err
	}

	j, err := genius.NewFromRawJSON(response)
	if err != nil {
		return nil, err
	}

	result := UserInfo{}
	result.Username = cast.ToString(j.Get("name"))
	result.Email = cast.ToString(j.Get("email"))
	return &result, nil
}
