FROM golang:1.24.1 AS build
WORKDIR /monitoring/
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o monitoring .

FROM scratch
COPY --from=build /monitoring/monitoring /monitoring/config.json /monitoring/
ENTRYPOINT [ "/monitoring/monitoring" ]