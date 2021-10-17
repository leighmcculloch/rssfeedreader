FROM golang as builder

WORKDIR /rssfeedreader
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o rssfeedreader

FROM gcr.io/distroless/static

COPY --from=builder /rssfeedreader/rssfeedreader /rssfeedreader

CMD ["/rssfeedreader"]
