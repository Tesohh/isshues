
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o isshues 

EXPOSE 2222

CMD ["/build/isshues"]
