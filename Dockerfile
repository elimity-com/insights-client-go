FROM golang
RUN \
    curl https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh && \
    go install github.com/jdeflander/goarrange@latest
