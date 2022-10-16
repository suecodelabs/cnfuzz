FROM golang:1.18 AS build

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /src
RUN make cnfuzz

FROM scratch AS final
COPY --from=build /src/dist/cnfuzz /
ENTRYPOINT ["/cnfuzz"]
