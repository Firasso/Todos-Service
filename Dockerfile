FROM golang:1.15.5-alpine as builder

ENV APP_DIR /go/src/github.com/Firasso/DemoGoApi
WORKDIR $APP_DIR

RUN apk add --no-cache bash make git tar curl

# Leverage go modules cache to speed up builds
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . $APP_DIR

RUN make build

# ---

FROM alpine

ENV APP_DIR /go/src/github.com/Firasso/DemoGoApi

COPY --from=builder $APP_DIR/build/Demo-app /bin/Demo-app
COPY --from=builder $APP_DIR/sql /bin/sql

WORKDIR /bin


EXPOSE 3000

ENTRYPOINT ["/bin/Demo-app"]
