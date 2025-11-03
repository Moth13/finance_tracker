FROM golang:1.25-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM scratch as runner
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env.sample app.env
EXPOSE 8080
CMD ["./main"]
