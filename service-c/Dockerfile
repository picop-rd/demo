# syntax=docker/dockerfile:1
FROM golang:1.19-bullseye AS builder

WORKDIR /go/src/github.com/hiroyaonoe/bcop/demo/service-c

RUN mkdir -p -m 0600 ~/.ssh \
	&& ssh-keyscan github.com >> ~/.ssh/known_hosts \
	&& git config --global url."git@github.com:".insteadOf "https://github.com/"
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /service-c main.go


FROM scratch

COPY --from=builder /service-c /bin/service-c
EXPOSE 80
ENTRYPOINT [ "/bin/service-c" ]
CMD [ "--port", "80", "--mysql-service", "user:password@tcp(proxy-c:80)/database?parseTime=true&collation=utf8mb4_bin", "--message", "This is service-c"]

