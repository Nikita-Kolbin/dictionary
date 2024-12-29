# TODO: Сделать двухэтапную сборку
#FROM golang:1.23 AS build
#
#COPY . /build/backend
#
#WORKDIR /build/backend/cmd
#
#RUN go build -o service ./main.go
#
#FROM alpine:3.20 AS runner
#
#COPY --from=build /build/backend/cmd /build
#
#ENV IN_DOCKER=1
#
#ENTRYPOINT ["/build/service"]

FROM golang:1.23

WORKDIR /backend

COPY . .

ENV IN_DOCKER=1

RUN go mod tidy
RUN go mod vendor

RUN go build -o build ./cmd/main.go

CMD ["./build"]
