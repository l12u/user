FROM golang:1.17.2-alpine AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# build the program
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o ./out/app ./cmd/userm

################
# run program
################

FROM alpine:3.10

COPY --from=builder /build/out/app .

CMD ["./app"]