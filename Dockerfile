FROM golang:1.12 as builder

WORKDIR /workdir
COPY . .

RUN CGO_ENABLED=0 go build -v -o rssfeedreader

FROM gcr.io/distroless/static

COPY --from=builder /workdir/rssfeedreader /rssfeedreader
COPY --from=builder /workdir/index.html /index.html

CMD ["/rssfeedreader"]
