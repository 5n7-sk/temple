<p align="center">
  <a href="https://github.com/skmatz/temple">
    <img src="./assets/images/banner.png" width="1000" alt="banner" />
  </a>
</p>

<p align="center">
  <a href="https://github.com/skmatz/temple/actions?query=workflow%3Abuild">
    <img
      src="https://github.com/skmatz/temple/workflows/build/badge.svg"
      alt="build"
    />
  </a>
  <a href="https://github.com/skmatz/temple/actions?query=workflow%3Arelease">
    <img
      src="https://github.com/skmatz/temple/workflows/release/badge.svg"
      alt="release"
    />
  </a>
  <a href="./LICENSE">
    <img
      src="https://img.shields.io/github/license/skmatz/temple"
      alt="license"
    />
  </a>
  <a href="./go.mod">
    <img
      src="https://img.shields.io/github/go-mod/go-version/skmatz/temple"
      alt="go version"
    />
  </a>
  <a href="https://github.com/skmatz/temple/releases/latest">
    <img
      src="https://img.shields.io/github/v/release/skmatz/temple"
      alt="release"
    />
  </a>
</p>

# Temple

**Temple** is a simple CUI for copying files at hand.

Have you ever wondered, "I need that file again, but where is it?"  
Temple brings the file to your hand quickly and accurately with CLI.

## Usage

```console
> temple

Search: █
Select
  ▸ ~/go/src/github.com/skmatz/temple/main.go

Name:           main.go
Path:           ~/go/src/github.com/skmatz/temple/main.go
Tags:           go, temple
Content:
   1 package main
   2
   3 import (
   4         "fmt"
   5         "io"
   6         "io/ioutil"
   7         "log"
   8         "net/http"
   9         "os"
  10         "os/user"
```

## Install

### Binary

Get binary from [releases](https://github.com/skmatz/temple/releases).  
If you already have [jq](https://github.com/stedolan/jq) and [fzf](https://github.com/junegunn/fzf) or [peco](https://github.com/peco/peco), you can download binary by running the following command.

```sh
curl -Ls https://api.github.com/repos/skmatz/temple/releases/latest | jq -r ".assets[].browser_download_url" | fzf | wget -i -
```

### Source

```sh
go get github.com/skmatz/temple
```
