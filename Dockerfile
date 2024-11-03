FROM docker.io/golang:1.23 AS builder

RUN apt-get update && apt-get install -y gcc

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN make install-prereq && \
    make build-linux

FROM scratch

COPY --from=builder /app/build/amd64/linux/note /app/note

ENTRYPOINT ["/app/note"]
