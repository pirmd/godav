.POSIX:

GO = go
LDFLAGS_STATIC = -linkmode 'external' -extldflags '-static'
LDFLAGS        = ${LDFLAGS_STATIC} -w -s

BIN        = godav
SRC_MOD    = go.mod
SRC       != ${GO} list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./...
SRC_TEST  != ${GO} list -f '{{range .TestGoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./...

all: ${BIN}

${BIN}: ${SRC} ${SRC_MOD} ${SRC_TEST}
	${GO} fmt ./...
	${GO} mod tidy -v
	${GO} mod verify
	${GO} test -vet=all ./...
	${GO} build -ldflags "${LDFLAGS}" -o $@
	grep -xq "$@" .gitignore || echo $@ >> .gitignore

audit: ${SRC} ${SRC_MOD}
	${GO} run honnef.co/go/tools/cmd/staticcheck@latest ./...
	${GO} run github.com/kisielk/errcheck@latest ./...
	${GO} run github.com/securego/gosec/v2/cmd/gosec@latest -quiet -exclude=G404 ./...
	${GO} run golang.org/x/vuln/cmd/govulncheck@latest ./...
.PHONY: audit

clean:
	${GO} clean
.PHONY: clean
