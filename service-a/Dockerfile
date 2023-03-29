# syntax=docker/dockerfile:1
FROM golang:1.19-bullseye AS builder

WORKDIR /go/src/github.com/picop-rd/demo/service-a

RUN mkdir -p -m 0600 ~/.ssh \
	&& ssh-keyscan github.com >> ~/.ssh/known_hosts \
	&& git config --global url."git@github.com:".insteadOf "https://github.com/"
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /service-a main.go


FROM scratch

COPY --from=builder /service-a /bin/service-a
EXPOSE 80
ENTRYPOINT [ "/bin/service-a" ]
CMD [ "--port", "80", "--child-service", "http://proxy-b:80", "--message", "This is service-a"]

