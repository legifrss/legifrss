# README

## Abstract

Everything is written in English, but you can [open issues](https://github.com/ldicarlo/legifrance-rss/issues/new) in French or English.

_Why this repository_: I was missing a way to see every change in the french law, but keep only some of them for later.

Legifrance provides a way to receive newsletters, but sadly they're not optimal.

Note that this project is at its ALPHA stage, it's more of a PoC for now.

## How to use

Just add the following feed to you RSS reader: `https://raw.githubusercontent.com/ldicarlo/legifrance-rss/master/feed/all.xml`

If you want a specific NATURE or AUTHOR, you can find them [here](https://github.com/ldicarlo/legifrance-rss/tree/master/feed). Just add the file name to `https://raw.githubusercontent.com/ldicarlo/legifrance-rss/master/feed/`

(Example: `https://raw.githubusercontent.com/ldicarlo/legifrance-rss/master/feed/all_Commission-nationale-du-débat-public.xml`)

It works in [Feedly](https://feedly.com).

## Nightly

If you want to check the last features you can add the following feed to you RSS reader: `https://github.com/ldicarlo/legifrance-rss/tree/nightly/feed`

## TODO

- Add valid RSS checker.
- Tech:
  - one xml file with all inside
  - 1h cache on requests
- Doc:
  - https://legifrss.github.io
  - all types
  - all authors
- Feats:
  - https://legifrss.org/latest => all
  - https://legifrss.org/latest?q=écologie => search all with term
  - https://legifrss.org/latest?q=écologie&author=Commission-nationale-du-débat-public&type=loi => search loi with term and author
