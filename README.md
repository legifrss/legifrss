# README

## Abstract

Everything is written in English, but you can [open issues](https://github.com/ldicarlo/legifrance-rss/issues/new) in French or English.

_Why this repository_: I was missing a way to see every change in the french law, but keep only some of them for later.

Legifrance provides a way to receive newsletters, but sadly they're not optimal.

Note that this project is at its ALPHA stage, it's more of a PoC for now.

## How to use

Read the doc at https://legifrss.github.io. If you need it translated in English please contact me.

## How to run

### Batch _(the part that fetch stuff from Legifrance)_

`go build -o bin/batch cmd/batch/batch.go && bin/batch`

### Server _(the part that serve RSS content)_

`go build -o bin/server cmd/server/server.go && bin/server`

### Technical stuff

- the server is based on https://github.com/gin-gonic/gin, and it's pretty good!
- the server uses a local `db.json` as a database (10 days of retention)
- the tests are not good, I know :/
- I use some personal ansible scripts to deploy

## TODO

- [X] Add valid RSS checker.
- Tech:
  - [X] one xml file with all inside
  - [X] 1h cache on requests
  - [X] SSL certificate and HTTPS enable (I mean it's 2021 wth)
- Doc:
  - [X] https://legifrss.github.io
  - [X] all types
  - [X] all authors
- Feats:
  - [X] https://legifrss.org/latest => all
  - [X] https://legifrss.org/latest?q=écologie => search all with term
  - [X] https://legifrss.org/latest?q=écologie&author=Commission-nationale-du-débat-public&type=loi => search law with term and author
- [X] Add twitter Bot: https://twitter.com/Legifrance1
