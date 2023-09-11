FROM golang as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN mkdir build && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/main cmd/mlinkubeapi/main.go


FROM scratch

COPY --from=builder /app/build/main /main

CMD ["/main"]
