# syntax=docker/dockerfile:1

FROM golang:1.23 AS build-stage

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-linkmode 'external' -extldflags '-static'" -o /simple-poker-game-infra


FROM alpine:latest AS run-stage

WORKDIR /
COPY --from=build-stage /simple-poker-game-infra /simple-poker-game-infra
COPY .env ./

EXPOSE 3000
CMD ["/simple-poker-game-infra"]
