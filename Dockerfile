FROM golang:1.18-alpine3.15 as build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ce-midtrans

FROM alpine:3.15

RUN apk add --no-cache ca-certificates

ENV K_SINK=
ENV SERVER_KEY=
ENV PORT=8080
WORKDIR /app

COPY --from=build /app/ce-midtrans .

EXPOSE 8080
CMD [ "./ce-midtrans" ]
