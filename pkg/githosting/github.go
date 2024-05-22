package githosting

import (
	"context"

	"github.com/jaronnie/genius"
	"github.com/jaronnie/restc"
	"github.com/spf13/cast"
)

type Github struct {
	Config Config
	Client restc.Interface
}

func (g *Github) GetUserInfo() (*UserInfo, error) {
	response, err := g.Client.Get().SubPath("/user").Do(context.Background()).RawResponse()
	if err != nil {
		return nil, err
	}
	j, err := genius.NewFromRawJSON(response)
	if err != nil {
		return nil, err
	}
	var result UserInfo

	result.Username = cast.ToString(j.Get("name"))
	result.Email = cast.ToString(j.Get("email"))
	return &result, nil
}
