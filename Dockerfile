FROM golang:1.23 as build

ENV CGO_ENABLED=0

ADD . /build
WORKDIR /build

RUN cd app && go build -o /build/ttag -ldflags "-s -w"

FROM golang:1.23

COPY --from=build /build/ttag /srv/app/ttag

WORKDIR /srv/app
CMD ["/srv/app/ttag"]
