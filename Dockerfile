FROM golang:1.17 AS build

WORKDIR /fuzzer

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /fuzzer
RUN make build

FROM scratch AS final
COPY --from=build /fuzzer/dist/cnfuzz /
EXPOSE 80
ENTRYPOINT ["/cnfuzz"]