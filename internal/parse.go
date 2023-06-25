package internal

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Parse struct {
	Url string
}

func NewParse(url string) *Parse {
	return &Parse{Url: url}
}

type Info struct {
	Protocol   string // 协议
	Address    string // git 地址
	User       string // 用户名
	Repository string // 仓库名
}

func (p *Parse) GetInfo() (*Info, error) {
	if strings.HasPrefix(p.Url, "git") {
		return getSSHInfo(p.Url)
	}
	if strings.HasPrefix(p.Url, "http") || strings.HasPrefix(p.Url, "https") {
		return getHTTPInfo(p.Url)
	}
	return nil, nil
}

func (p *Parse) Replace() (string, error) {
	i, err := p.GetInfo()
	if err != nil {
		return "", err
	}

	if strings.Contains(i.Address, "github") {
		// github, get GITHUB_TOKEN env
		s := os.Getenv("GITHUB_TOKEN")
		if s == "" {
			return "", errors.New("empty GITHUB_TOKEN")
		}
		return fmt.Sprintf("%s://%s@github.com/%s/%s.git", i.Protocol, s, i.User, i.Repository), nil
	}

	return "", errors.New("unknown error")
}

func getSSHInfo(u string) (*Info, error) {
	x, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	host := x.Host
	userName := x.Path[1:strings.Index(x.Path, "/")]
	repoName := x.Path[strings.Index(x.Path, "/")+1 : len(x.Path)-len(".git")]
	return &Info{
		Protocol:   x.Scheme,
		Address:    host,
		User:       userName,
		Repository: repoName,
	}, nil
}

func getHTTPInfo(u string) (*Info, error) {
	x, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	host := x.Host
	userName := strings.Split(x.Path[1:], "/")[0]
	repoName := strings.TrimSuffix(strings.Split(x.Path[1:], "/")[1], ".git")
	return &Info{
		Protocol:   x.Scheme,
		Address:    host,
		User:       userName,
		Repository: repoName,
	}, nil
}
