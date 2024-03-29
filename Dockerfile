FROM golang

WORKDIR /rssfeedreader
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o rssfeedreader

FROM gcr.io/distroless/static

COPY --from=0 /rssfeedreader/rssfeedreader /rssfeedreader

CMD ["/rssfeedreader"]
