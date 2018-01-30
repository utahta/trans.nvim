# trans

[![Go Report Card](https://goreportcard.com/badge/github.com/utahta/trans)](https://goreportcard.com/report/github.com/utahta/trans)

A Google Translate CLI written in Go.

## Required

You must first set up authentication by creating a service account.

The service account documentation can be found [here](https://cloud.google.com/iam/docs/creating-managing-service-accounts).

After that, set the environment variable `GOOGLE_APPLICATION_CREDENTIALS` to the file path of the JSON file that contains your service account key.

## Note

Need a little bit of money to use google translate API.

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

