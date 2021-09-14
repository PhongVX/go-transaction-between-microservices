FROM golang:1.16.5 as builder
ARG SOURCE_LOCATION=/
WORKDIR ${SOURCE_LOCATION}
COPY go.* ./
RUN go env -w GOFLAGS=-mod=mod \
   && go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
ARG SOURCE_LOCATION=/
RUN apk --no-cache add curl
EXPOSE 8082
WORKDIR /root/
COPY --from=builder ${SOURCE_LOCATION}/app ./app

CMD ["./app"]  
