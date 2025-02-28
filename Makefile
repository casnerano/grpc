GO_COVER_EXCLUDE_PATTERN = "\/mocks\/"

.PHONY: test
test:
	go test ./...

.PHONY: test-cover
test-cover:
	go test -count=1 -cover -coverprofile=cover.temp.out -covermode=atomic ./...
	grep -vE ${GO_COVER_EXCLUDE_PATTERN} cover.temp.out > cover.out && rm cover.temp.out
	go tool cover -func=cover.out

.PHONY: test-cover-with-html
test-cover-with-html: test-cover
	go tool cover -html=cover.out