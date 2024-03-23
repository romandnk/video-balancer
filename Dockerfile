FROM golang:1.22.0 as build

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/app ./cmd/app

FROM scratch

WORKDIR /app

COPY --from=build /app/bin/app ./bin/
COPY --from=build /app/config/config.yaml ./config/config.yaml

ENTRYPOINT ["./bin/app"]