FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod vendor

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM scratch

WORKDIR /app
COPY --from=builder /app/main /main
COPY --from=builder /app/config.yml /app/config.yml 
COPY --from=builder /app/programs.json /app/programs.json 

CMD ["/main"]
