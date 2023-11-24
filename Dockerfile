FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o /core.com

VOLUME [ "/mydata" ]
# RUN mkdir -p /app/data

CMD [ "/core.com" ]