WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN sqlc generate

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN goose up

RUN go build -o isshues 

EXPOSE 2222

CMD ["/build/isshues"]
