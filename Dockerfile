FROM golang:alpine AS server-build

RUN apk add --no-cache git

RUN adduser -D -u 10000 jan
RUN mkdir /server && chown jan /server/
USER jan

WORKDIR /server/
COPY . /server/

RUN CGO_ENABLED=0 go build -o ./go-micro .

FROM alpine

WORKDIR /
COPY --from=server-build /server/go-micro .

ENV SERVICE_ADDR=:8080
EXPOSE 8080
CMD ["/go-micro"]