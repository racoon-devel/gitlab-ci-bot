FROM golang as builder
WORKDIR /go/src/github.com/racoon-devel/gitlab-ci-bot
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gitlab-ci-bot -a -installsuffix cgo ./app/main.go
FROM alpine:latest
RUN mkdir /app && mkdir /etc/gitlab-ci-bot
WORKDIR /app
COPY --from=builder /go/src/github.com/racoon-devel/gitlab-ci-bot/gitlab-ci-bot .
COPY --from=builder /go/src/github.com/racoon-devel/gitlab-ci-bot/configs/config.toml /etc/gitlab-ci-bot/config.toml
EXPOSE 8080/tcp
CMD ["./gitlab-ci-bot"]