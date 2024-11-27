FROM golang:alpine AS build

# We need tzdata for the timezone information and the
# ca-certificates for ssl cert verification
RUN apk --no-cache add tzdata ca-certificates

WORKDIR /app

COPY main.go ./

COPY go.* ./

RUN go build -a -tags netgo -ldflags '-w' -v -o main .

FROM scratch

WORKDIR /app

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/main .

CMD [ "/app/main" ]