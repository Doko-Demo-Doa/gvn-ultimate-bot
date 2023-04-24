FROM golang:1.20.3

ENV LOG_LEVEL=info
ENV GIN_MODE release

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN ls -l

RUN go build -o gvn-ultimate-bot

EXPOSE 3000

CMD [ "./gvn-ultimate-bot" ]
