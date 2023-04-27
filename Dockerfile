FROM golang:1.19 as builder

WORKDIR /go/src

COPY ./go.mod ./go.sum ./

ENV GO111MODULE=on

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./app ./main.go

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /go/src/app /go/src/app

ENV GIN_MODE=release

CMD ["/go/src/app"]
