[![build](https://github.com/linkdata/jawsdemo/actions/workflows/build.yml/badge.svg)](https://github.com/linkdata/jawsdemo/actions/workflows/build.yml)
[![goreport](https://goreportcard.com/badge/github.com/linkdata/jawsdemo)](https://goreportcard.com/report/github.com/linkdata/jawsdemo)

# jawsdemo

Simple demo application for [JaWS](https://github.com/linkdata/jaws), with extensive comments on how to use it.

Maybe start by having a look at the comments in `assets/ui/index.html`.

## Use

Generate the version file, then run the application:

`go generate ./...`

`go run .`

Open multiple web browser pages to `http://localhost:8080/`. Play around with the controls on the pages. As you alter the values in one web page, the others match the changes in real time.

To use a different port, pass an explicit address:

`go run . -address localhost:8081`
