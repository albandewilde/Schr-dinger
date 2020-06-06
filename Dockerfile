FROM golang:1.14 as builder

WORKDIR /go/src/s_cat
COPY . .

RUN CGO_ENABLED=0 go build -o /bin/s_cat


FROM alpine

WORKDIR /bin/s_cat

COPY --from=builder /bin/s_cat .
COPY secrets.json .

CMD ["./s_cat"]
