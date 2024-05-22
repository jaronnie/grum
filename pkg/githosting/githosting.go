package githosting

import (
	"fmt"
	"net/http"

	"github.com/jaronnie/restc"
	"github.com/pkg/errors"
)

const (
	GITHUB = "github"
	GITLAB = "gitlab"
	GITEA  = "gitea"
)

const (
	AuthorizationPrefixKey = "Bearer"
)

type Config struct {
	Type  string
	Url   string
	Token string
}

type UserInfo struct {
	Username string
	Email    string
}

type Interface interface {
	GetUserInfo() (*UserInfo, error)
}

func New(config Config) (Interface, error) {
	var restClient restc.Interface
	var err error
	headers := make(http.Header)

	switch config.Type {
	case GITHUB:
		headers.Set("Authorization", fmt.Sprintf("%s %s", AuthorizationPrefixKey, config.Token))
		restClient, err = restc.New(restc.WithUrl("https://api.github.com"), restc.WithHeaders(headers))
		return &Github{Config: config, Client: restClient}, err
	case GITLAB:
		headers.Set("PRIVATE-TOKEN", config.Token)
		restClient, err = restc.New(restc.WithUrl(config.Url), restc.WithHeaders(headers))
		return &Gitlab{Config: config, Client: restClient}, err
	case GITEA:
		headers.Set("Authorization", fmt.Sprintf("%s %s", AuthorizationPrefixKey, config.Token))
		restClient, err = restc.New(restc.WithUrl(config.Url), restc.WithHeaders(headers))
		return &Gitea{Config: config, Client: restClient}, err
	default:
		return nil, errors.Errorf("not support %s type", config.Type)
	}
}
