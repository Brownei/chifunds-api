FROM golang:alpine

RUN go install github.com/air-verse/air@latest

WORKDIR /chifunds

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8000

CMD ["air"]
