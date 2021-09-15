FROM golang:1.16 AS builder

WORKDIR /go/builder
# ADD go.mod go.sum ./
ADD go.mod ./
RUN go mod download -x

ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o client cmd/client/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o server cmd/server/main.go

# final stage
FROM scratch

COPY --from=builder /go/builder/client /client
COPY --from=builder /go/builder/server /server
