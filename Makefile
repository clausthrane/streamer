MOCKGEN=mockgen
MOCKDIR=.mocks
PACKAGES	?= $(shell go list ./...)
FILES		?= $(shell find . -type f -name '*.go' -not -path "./.mocks/*")

tools:
	go get -u github.com/golang/mock/mockgen
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint
	go get -u go.uber.org/goleak

test: mocks
	go test -race ./...

cover: mocks
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o=cover.html

fmt:
	go fmt ./...
	goimports -w $(FILES)

lint:
	golint $(PACKAGES)

vet:
	go vet ./...

mocks:
	$(MOCKGEN) -destination $(MOCKDIR)/interface.go -source interface.go Connector -destination $(MOCKDIR)
	$(MOCKGEN) -destination $(MOCKDIR)/source/stream.go -source source/stream.go Stream -destination $(MOCKDIR)

clean-mocks:
	rm -rf $(MOCKDIR) && mkdir $(MOCKDIR)
	make mocks