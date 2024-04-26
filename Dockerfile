FROM golang:latest as build

WORKDIR /src/
COPY . /src/
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /src/graphql-go-example .

FROM gcr.io/distroless/base-debian11
COPY --from=build /src/graphql-go-example /usr/bin/graphql-go-example

EXPOSE 8080
ENTRYPOINT ["/usr/bin/graphql-go-example", "-address=0.0.0.0"]