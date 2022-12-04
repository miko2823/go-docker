FROM golang:1.19.3-alpine

RUN apk update && apk add git

RUN go install golang.org/x/tools/cmd/goimports@latest

WORKDIR /go/src/app

# gocode-gomod
# RUN go get -x -d github.com/stamblerre/gocode \
#     && go build -o gocode-gomod github.com/stamblerre/gocode \
#     && mv gocode-gomod $GOPATH/bin/

# [Optional] Uncomment the next line to use go get to install anything else you need
# RUN go get -x <your-dependency-or-tool>

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1
