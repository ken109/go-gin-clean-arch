FROM golang:1.23 as dev

WORKDIR /go/src

RUN go install github.com/air-verse/air@v1.61.5
RUN go install github.com/maoueh/zap-pretty/cmd/zap-pretty@v0.3.1

COPY ./go.mod ./go.sum ./

ENV GO111MODULE=on

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./app ./cmd/api/main.go

FROM gcr.io/distroless/static:nonroot

COPY --from=dev /go/src/app /go/src/app

ENV GIN_MODE=release

CMD ["/go/src/app"]
