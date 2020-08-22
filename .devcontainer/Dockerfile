FROM golang:latest

RUN curl -fsSL https://get.docker.com -o get-docker.sh
RUN sh get-docker.sh
RUN rm get-docker.sh

RUN go get -v golang.org/x/tools/cmd/goimports \
              golang.org/x/tools/gopls \
              golang.org/x/lint/golint \
              github.com/ramya-rao-a/go-outline \
              github.com/mdempsky/gocode \
              github.com/uudashr/gopkgs/v2/cmd/gopkgs \
              github.com/go-delve/delve/cmd/dlv \
              github.com/rogpeppe/godef

EXPOSE 3000
CMD [ "echo", "\"This image is intended only for use inside of Visual Studio Code.\"" ]
