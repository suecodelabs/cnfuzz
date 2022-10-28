FROM golang:1.18 AS build

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /src
RUN make restlerwrapper

FROM mcr.microsoft.com/restlerfuzzer/restler:v8.6.0 as final
COPY --from=build /src/dist/restlerwrapper /
COPY src/cmd/restlerwrapper/auth.py /scripts/auth.py

