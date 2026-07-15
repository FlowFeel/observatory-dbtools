FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download 2>/dev/null || true
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /dbtools ./cmd/dbtools

FROM alpine:3.20
RUN apk add --no-cache mysql-client bash curl
COPY --from=build /dbtools /usr/local/bin/dbtools
COPY migrations/ /migrations/
COPY baselines/ /baselines/
COPY Taskfile.yml /work/Taskfile.yml

# Install go-task
RUN curl -sL https://taskfile.dev/install.sh | sh -s -- -b /usr/local/bin

WORKDIR /work
ENTRYPOINT ["task"]
