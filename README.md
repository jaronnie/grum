# grum

git reomote url modify, for example:
It will modify your_repository/.git/config remote.origin url value.

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
go install github.com/jaronnie/grum
```

## usage

```shell
cd /path/to/your_repository
grum
# grum --type gitlab 使用公司内网的 gitlab
# grum --type gitlab --insecure # 使用 http
```


