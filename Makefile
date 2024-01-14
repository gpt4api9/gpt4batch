GOPATH:=$(shell go env GOPATH)

.PHONY: init
beta:
	@goreleaser check
	@goreleaser --snapshot --skip-publish --rm-dist

willRelease:
	@goreleaser check
	@goreleaser release --skip=publish

release:
	@goreleaser check
	@goreleaser release --snapshot --rm-dist
