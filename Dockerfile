FROM golang:1.18 AS build

WORKDIR /fuzzer

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /fuzzer
RUN make build

FROM scratch AS final
COPY --from=build /fuzzer/dist/cnfuzz /
EXPOSE 8080
ENTRYPOINT ["/cnfuzz"]
