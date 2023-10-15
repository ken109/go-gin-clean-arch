FROM golang:1.21 as dev

WORKDIR /go/src

RUN go install github.com/cosmtrek/air@latest

COPY ./go.mod ./go.sum ./

ENV GO111MODULE=on

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./app ./cmd/api/main.go

FROM gcr.io/distroless/static:nonroot

COPY --from=dev /go/src/app /go/src/app

ENV GIN_MODE=release

CMD ["/go/src/app"]
