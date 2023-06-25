# grum

git reomote url modify, for example:
it will modify your_repository/.git/config remote.origin url value

[remote "origin"]
fetch = +refs/heads/*:refs/remotes/origin/*
url = https://github.com/jaronnie/grum.git

to 

[remote "origin"]
    fetch = +refs/heads/*:refs/remotes/origin/*
    url = https://your_github_token@github.com/jaronnie/grum.git


support ssh or http protocol

## install

```shell
go install github.com/jaronnie/grum
```

## usage

```shell
cd /path/to/your_repository
grum
```