FROM golang as build-go
WORKDIR /lmApi
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/lmApi .

FROM alpine:latest
RUN addgroup -S lmApi && adduser -S lmApi -G lmApi
USER lmApi
WORKDIR /home/lmApi
COPY --from=build-go /bin/lmApi ./
EXPOSE 3000
ENTRYPOINT ["./lmApi"]