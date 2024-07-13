FROM golang:1.22.5

WORKDIR /app

COPY go.mod .
COPY main.go .
COPY static/quran-final.json ./static/quran-final.json

RUN go get

RUN go build -o main .
COPY static/quran-final.json ./static/quran-final.json

ENV GIN_MODE=release
ENTRYPOINT ["/app/main"]