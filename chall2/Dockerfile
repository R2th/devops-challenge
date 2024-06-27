FROM golang:alpine
WORKDIR /app
COPY src/config.env .

COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/. .
RUN go build -o app src/main.go

EXPOSE 9915

CMD ["app"]