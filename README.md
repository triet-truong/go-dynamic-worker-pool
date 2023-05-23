# Go dynamic worker pool

## Introduction

A Go worker pool with the ability of dynamical size modification (add or remove workers) using external configuration file (concurrency.txt).

## User Guide

Run the program:

```shell
go run main.go
```

A file named `concurrency.txt` will be created and contain the initial number of workers (default: `3`).
Update the number in this file and watch the program log:

- If the new number is higher than the current ones, new workers will be added.
- If the new number is lower than the current ones, some workers will be removed.
- If the new number is 0, all workers will be removed and subsequently program will exit.

<img width="914" alt="image" src="https://user-images.githubusercontent.com/19305944/211184964-3d95b980-8837-4022-abaa-31d532a2dcbd.png">

Try it yourself!

## Integration test coverage

Run program with coverage files generation:

```shell
mkdir covdatafiles
GOCOVERDIR=covdatafiles go run -cover *.go
```

Merge duplicated coverage files (if necessary):

```shell
mkdir merged
go tool covdata merge -i=covdatafiles -o=merged
```

Show coverage overall percentage:

```shell
go tool covdata percent -i=covdatafiles
```

Convert to human-readable coverage file:

```shell
go tool covdata textfmt -i=covdatafiles -o=cov.txt
```

Show as console output or HTML :

```shell
go tool cover -func=cov.txt
go tool cover -html=cov.txt
```
