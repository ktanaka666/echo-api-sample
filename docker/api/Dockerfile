FROM golang:1.11-alpine as build

RUN apk add git

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o bin/app

FROM build as production
COPY --from=build /go/src/app/bin/app bin/app

CMD ["bin/app"]
