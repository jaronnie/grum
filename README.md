# grum

git remote url modify, for example:
It will modify your_repository/.git/config remote.origin url value and user.name, user.email.

It is no longer necessary to configure cumbersome SSH Key, and the configuration can be completed with one click through the access token.

```shell
[remote "origin"]
    fetch = +refs/heads/*:refs/remotes/origin/*
    url = https://github.com/jaronnie/grum.git
```

to

```shell
[remote "origin"]
    fetch = +refs/heads/*:refs/remotes/origin/*
    url = https://your_github_token@github.com/jaronnie/grum.git
```

support ssh or http/https protocol.

## install

```shell
go install github.com/jaronnie/grum@latest
```

## usage

```shell
cd /path/to/your_repository
grum
# grum --type gitlab 使用公司内网的 gitlab
# grum --type gitlab --insecure # 使用 http
```

```shell
# default github
grum clone git@git.hyperchain.cn:niejian/sc.git --type gitlab
# http insecure
grum clone http://git.hyperchain.cn/niejian/sc.git --type gitlab --insecure
```

**Required environment variables, judged according to the type**

| type   | env name                                               |
| ------ | ------------------------------------------------------ |
| github | [GITHUB_TOKEN](https://github.com/settings/tokens/new) |
| gitlab | GITLAB_TOKEN                                           |

