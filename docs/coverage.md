## Code Coverage Reports for Golang Applications ##

This document covers how to generate HTML Code Coverage Reports for
Golang Applications.

#### Generate a test executable which calls your main()

```sh
$ go test -c -covermode=count -coverpkg ./...
```

#### Run the generated application to produce a new coverage report

```sh
$ ./sms.test -test.run "^TestMain$" -test.coverprofile=coverage.cov
```

#### Run your unit tests to produce their coverage report

```sh
$ go test -test.covermode=count -test.coverprofile=unit.out ./...
```

#### Merge the two coverage Reports

```sh
$ go get github.com/wadey/gocovmerge
$ gocovmerge unit.out coverage.cov > all.out
```

#### Generate HTML Report

```sh
$ go tool cover -html all.out -o coverage.html
```

#### Generate Function Report

```sh
$ go tool cover -func all.out
```