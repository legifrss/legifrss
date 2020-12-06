# README

## Abstract

Everything is written in English, but you can open issues in French.

_Why this repository_: I was missing a way to see every change in the french law, but keep only some of them for later.

Legifrance provides a way to receive newsletters, but sadly they're not optimal.

Note that this project is at its ALPHA stage, it's more of a PoC for now.

## How to use

Just add the following feed to you RSS reader: `https://raw.githubusercontent.com/ldicarlo/legifrance-rss/master/feed/feed.xml`

It works in [Feedly](https://feedly.com).

## Expected Features

I will soon add these:

- Fix: do not update date when contents are identical
- Feature: add category field to be able to sort by Nature and Author (_i.e._ Minister)
- Tests: I should add unit tests.

## Technical (WIP)

### Misc

Init using [go article](https://blog.golang.org/using-go-modules)
Highly based on [hnrss](https://github.com/hnrss)

### Auth

```bash
curl -k -v -H 'content-type: application/x-www-form-urlencoded' -d 'grant_type=client_credentials&client_id=CLIENT_ID&client_secret=CLIENT_SECRET&scope=openid'-X POST 'https://oauth.aife.economie.gouv.fr/api/oauth/token/api/oauth/token'
```

### To consult Last Jorfs

```bash
curl -is -H 'Authorization: Bearer [...]' -X POST 'https://sandbox-api.aife.economie.gouv.fr/dila/legifrance-beta/lf-engine-app/consult/lastNJo' -d '{"nbElement":5}'
```

Returns for example:

```json
{
  "executionTime": 12,
  "containers": [
    {
      "id": "JORFCONT000042612872",
      "ancienId": null,
      "idEli": "/eli/jo/2020/12/5/0294",
      "origine": "JORF",
      "url": "conteneur/JORF/CONT/00/00/42/61/28/JORFCONT000042612872.xml",
      "nature": "JO",
      "titre": "JORF n°0294 du 5 décembre 2020",
      "etat": null,
      "num": "0294",
      "numero": null,
      "datePubli": 1607126400000,
      "relevantDate": null,
      "refInjection": "MD-20201205_011008_743_JOQUOT",
      "idTechInjection": "index_DATA_JORF_CONT_JORFCONT000042612872"
    },
    {
      "id": "JORFCONT000042607089",
      "ancienId": null,
      "idEli": "/eli/jo/2020/12/4/0293",
      "origine": "JORF",
      "url": "conteneur/JORF/CONT/00/00/42/60/70/JORFCONT000042607089.xml",
      "nature": "JO",
      "titre": "JORF n°0293 du 4 décembre 2020",
      "etat": null,
      "num": "0293",
      "numero": null,
      "datePubli": 1607040000000,
      "relevantDate": null,
      "refInjection": "MD-20201204_185514_981_BDJQUOT",
      "idTechInjection": "index_DATA_JORF_CONT_JORFCONT000042607089"
    },
    {
      "id": "JORFCONT000042601748",
      "ancienId": null,
      "idEli": "/eli/jo/2020/12/3/0292",
      "origine": "JORF",
      "url": "conteneur/JORF/CONT/00/00/42/60/17/JORFCONT000042601748.xml",
      "nature": "JO",
      "titre": "JORF n°0292 du 3 décembre 2020",
      "etat": null,
      "num": "0292",
      "numero": null,
      "datePubli": 1606953600000,
      "relevantDate": null,
      "refInjection": "MD-20201204_105509_210_BDJQUOT",
      "idTechInjection": "index_DATA_JORF_CONT_JORFCONT000042601748"
    },
    {
      "id": "JORFCONT000042592258",
      "ancienId": null,
      "idEli": "/eli/jo/2020/12/2/0291",
      "origine": "JORF",
      "url": "conteneur/JORF/CONT/00/00/42/59/22/JORFCONT000042592258.xml",
      "nature": "JO",
      "titre": "JORF n°0291 du 2 décembre 2020",
      "etat": null,
      "num": "0291",
      "numero": null,
      "datePubli": 1606867200000,
      "relevantDate": null,
      "refInjection": "MD-20201203_205511_354_BDJQUOT",
      "idTechInjection": "index_DATA_JORF_CONT_JORFCONT000042592258"
    },
    {
      "id": "JORFCONT000042580121",
      "ancienId": null,
      "idEli": "/eli/jo/2020/12/1/0290",
      "origine": "JORF",
      "url": "conteneur/JORF/CONT/00/00/42/58/01/JORFCONT000042580121.xml",
      "nature": "JO",
      "titre": "JORF n°0290 du 1 décembre 2020",
      "etat": null,
      "num": "0290",
      "numero": null,
      "datePubli": 1606780800000,
      "relevantDate": null,
      "refInjection": "MD-20201204_185514_981_BDJQUOT",
      "idTechInjection": "index_DATA_JORF_CONT_JORFCONT000042580121"
    }
  ],
  "totalNbResult": 0
}

```

Then fetching with:

```bash
curl  -H "Authorization: Bearer $JWT" -H "Content-Type: application/json" -X POST 'https://api.aife.economie.gouv.fr/dila/legifrance-beta/lf-engine-app/consult/jorfCont' -d '{"textCid":"JORFCONT000042580121"}' -w '\n' | jq
```
