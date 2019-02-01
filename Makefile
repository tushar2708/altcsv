
GLIDEBIN=${GOPATH}/bin/glide
GOFILES=$(find . -path ./vendor -prune -o -name '*.go' -print)

$(GLIDEBIN):
	go get github.com/Masterminds/glide

install: $(GLIDEBIN)
	${GOPATH}/bin/glide install

test: install
	go build
	go test -v ./
	# gofmt -l ${GOFILES} | read && echo "gofmt failures" && gofmt -d ${GOFILES} && exit 1 || true

clean:
	rm -rf build