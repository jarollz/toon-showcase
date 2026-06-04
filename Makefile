SHELL := /bin/bash

PKGS := ./...
COVERPKG := ./...
COVERPROFILE := .coverage.out
COVERMODE := atomic
COVER_MIN := 90.0

.PHONY: help test test-race vet cover cover-check clean ci

help:
	@printf "%s\n" \
	"make test         - run unit tests" \
	"make test-race    - run tests with race detector" \
	"make vet          - go vet" \
	"make cover        - run coverpkg coverage report" \
	"make cover-check  - fail if total coverage < COVER_MIN" \
	"make ci           - vet + test-race + cover-check" \
	"make clean        - remove coverage artifacts"

test:
	go test $(PKGS)

test-race:
	go test -race $(PKGS)

vet:
	go vet $(PKGS)

cover:
	go test -coverpkg=$(COVERPKG) -covermode=$(COVERMODE) -coverprofile=$(COVERPROFILE) $(PKGS)
	go tool cover -func=$(COVERPROFILE)

cover-check: cover
	@total=$$(go tool cover -func=$(COVERPROFILE) | awk '/^total:/ {gsub("%","",$$3); print $$3}'); \
	echo "total coverage: $$total% (min: $(COVER_MIN)%)"; \
	awk -v t="$$total" -v m="$(COVER_MIN)" 'BEGIN {exit (t+0 >= m+0) ? 0 : 1}'

clean:
	rm -f $(COVERPROFILE)

ci: vet test-race cover-check
