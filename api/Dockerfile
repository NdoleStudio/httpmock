FROM golang:latest as builder

ARG GIT_COMMIT
ENV GIT_COMMIT=$GIT_COMMIT

WORKDIR /httpmock

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .


RUN go install github.com/swaggo/swag/cmd/swag@v1.8.7
RUN swag init

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$GIT_COMMIT" -o /bin/httpmock .

FROM alpine:latest

RUN addgroup -S httpmock && adduser -S httpmock -G httpmock

USER httpmock
WORKDIR /home/httpmock

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /zoneinfo.zip
COPY --from=builder /bin/httpmock ./

ENV ZONEINFO=/zoneinfo.zip

EXPOSE 8000

ENTRYPOINT ["./httpmock", "--dotenv=false"]
