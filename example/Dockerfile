FROM golang:1.18 AS build
# RUN apt-get -y update
RUN go install github.com/swaggo/swag/cmd/swag@v1.7.9
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN make swag
# RUN CGO_ENABLED=0 GOOS=linux go build -o todoapi main.go
RUN mkdir -p dist
RUN make build
RUN ls dist

FROM alpine as final
COPY --from=build /app/dist .
CMD ["./todo-api"]
