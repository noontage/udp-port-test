FROM golang:1.18-buster as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go mod download && go build -o /go/bin/app /go/src/app/cmd/server/main.go

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/app /
EXPOSE 8081 39998
CMD ["/app"]