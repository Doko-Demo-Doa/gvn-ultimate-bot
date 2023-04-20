FROM golang:1.20.1

ENV LOG_LEVEL=info
ENV ENV=production
ENV GIN_MODE release

WORKDIR /

COPY go.mod ./
RUN go mod download

RUN go build -o gvn-ultimate-bot

EXPOSE 3000

CMD [ "./gvn-ultimate-bot" ]
