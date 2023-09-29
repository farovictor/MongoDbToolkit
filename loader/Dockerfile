FROM golang:1.20-bullseye as builder

ARG GIT_COMMIT
ARG BUILD_TIME
ARG VERSION

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags "-X 'github.com/farovictor/MongoDbLoader/src/cmd.GitCommit=${GIT_COMMIT}' -X 'github.com/farovictor/MongoDbLoader/src/cmd.BuildTime=${BUILD_TIME}' -X 'github.com/farovictor/MongoDbLoader/src/cmd.Version=${VERSION}' -s -w" -o mongoload -a -installsuffix cgo ./loader/src && \
    chmod a+x mongoload

FROM alpine:latest

COPY --from=builder /app/mongoload /usr/bin/mongoload

ENTRYPOINT [ "mongoload" ]