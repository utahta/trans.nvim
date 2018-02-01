# trans

[![Go Report Card](https://goreportcard.com/badge/github.com/utahta/trans)](https://goreportcard.com/report/github.com/utahta/trans)

Google Translate CLI written in Go.

## Required

You must first set up authentication by creating a API Key or service account of GCP.

### API Key

The API Key documentation can be found [here](https://cloud.google.com/translate/docs/auth#using_an_api_key).

Set the environment variable `TRANS_API_KEY` to the API Key.

### Service Account

The service account documentation can be found [here](https://cloud.google.com/iam/docs/creating-managing-service-accounts).

Set the environment variable `GOOGLE_APPLICATION_CREDENTIALS` to the file path of the JSON file that contains your service account key.

## Note

You need a little bit of money to use google translate API.  
e.g. it costs $0.06 for 2889 characters.

## Installation

```sh
$ go get -u github.com/utahta/trans/cmd/trans
```

## Usage

```sh
$ trans -t en こんにちは
Hello
```
```sh
$ trans -t ja this is a pen
これはペンです
```

