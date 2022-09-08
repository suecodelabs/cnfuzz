FROM golang:1.18 AS build

WORKDIR /wrapper

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /wrapper
RUN make restlerwrapper

FROM mcr.microsoft.com/restlerfuzzer/restler:v8.6.0 as final
COPY --from=build /wrapper/dist/restlerwrapper /
EXPOSE 8080
ENTRYPOINT ["/restlerwrapper"]
