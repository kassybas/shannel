

install:
	go install github.com/fatih/errwrap@latest
	go get -v -u github.com/go-critic/go-critic/cmd/gocritic

gocritic:
	gocritic check -disable=appendAssign -enable="#style,#performance,#opinionated,#security,#diagnostic" ./...

errwrap:
	errwrap ./...

golint:
	golint -set_exit_status ./...

staticcheck:
	staticcheck ./...

static-analysis: errwrap golint staticcheck gocritic
	