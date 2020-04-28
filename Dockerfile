###################
##  build stage  ##
###################
FROM golang:1.13.0-alpine as builder
WORKDIR /profileserver-golang-kubernetes
COPY . .
RUN go build -v -o profileserver-golang-kubernetes

##################
##  exec stage  ##
##################
FROM alpine:3.10.2
WORKDIR /app
COPY ./configs/config.json.default ./configs/config.json
COPY --from=builder /profileserver-golang-kubernetes /app/
CMD ["./profileserver-golang-kubernetes"]
