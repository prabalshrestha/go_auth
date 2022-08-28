FROM golang:1.17.2-alpine3.14
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build main.go
CMD [ "./main" ]