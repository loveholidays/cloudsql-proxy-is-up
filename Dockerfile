FROM golang:1.15 as golang
WORKDIR /go/src/github.com/loveholidays/cloudsql-proxy-is-up/
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cloudsql-proxy-is-up .


FROM gcr.io/distroless/base-debian10

COPY --from=golang /go/src/github.com/loveholidays/cloudsql-proxy-is-up/cloudsql-proxy-is-up .
ENTRYPOINT ["/cloudsql-proxy-is-up"]
