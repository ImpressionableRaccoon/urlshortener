version=$(shell git describe --always --long --dirty)
date=$(shell TZ=UTC date)
commit=$(shell git log -1 --pretty=format:"%H")

all:
	go build -o shortener -ldflags '\
-X "main.buildVersion=${version}" -X "main.buildDate=${date}" -X "main.buildCommit=${commit}"' cmd/shortener/main.go

staticlint:
	go build -o staticlint cmd/staticlint/main.go