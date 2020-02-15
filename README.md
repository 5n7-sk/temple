<p align="center">
  <a href="https://github.com/skmatz/temple">
    <img src="./assets/images/banner.png" width="1000" alt="banner" />
  </a>
</p>

<p align="center">
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

Path:           ~/go/src/github.com/skmatz/temple/main.go
Tags:           go, temple
Content:
  package main

  import (
          "fmt"
          "io"
          "io/ioutil"
          "log"
          "net/http"
          "os"
          "os/user"
```

## Install

### Binary

Get binary from [releases](https://github.com/skmatz/temple/releases).

### Source

```sh
go get github.com/skmatz/temple
```
