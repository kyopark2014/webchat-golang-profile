###################
##  build stage  ##
###################
FROM golang:1.13.0-alpine as builder
WORKDIR /webchat-golang-profile
COPY . .
RUN go build -v -o webchat-golang-profile

##################
##  exec stage  ##
##################
FROM alpine:3.10.2
WORKDIR /app
COPY ./configs/config.json.default ./configs/config.json
COPY --from=builder /webchat-golang-profile /app/
CMD ["./webchat-golang-profile"]
