FROM golang AS compiling_stage
RUN mkdir -p /go/src/CityTransport
WORKDIR /go/src/CityTransport
ADD main.go .
ADD go.mod .
RUN go build -o transport_simulator .

FROM alpine:latest
LABEL version="1.0.0"
LABEL maintainer="DENIS KOZLOV"
WORKDIR /root/
COPY --from=compiling_stage /go/src/CityTransport/transport_simulator .
ENTRYPOINT ["./transport_simulator"]