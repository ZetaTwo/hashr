all: hashr query_tool

hashr:
	env GOOS=linux GOARCH=amd64 go build hashr.go

query_tool:
	env GOOS=linux GOARCH=amd64 go build tool/lookup.go

clean:
	rm -f hashr lookup

.PHONY: all clean
