FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/tlehman/goga.me

ENV USER tlehman
ENV HTTP_ADDR 8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET Dq7K0yUtpbH7bBrY

# Replace this with actual PostgreSQL DSN.
ENV DSN postgres://tlehman@localhost:5432/goga.me?sslmode=disable

WORKDIR /go/src/github.com/tlehman/goga.me

RUN godep go build
RUN ./goga.me

EXPOSE 8888