FROM golang:1.18

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY / ./

RUN go build -o ./strava-api

EXPOSE 8080

CMD [ "./strava-api" ]