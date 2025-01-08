FROM golang:1.22-alpine

ENV USER=cert400
ENV GROUPNAME=$USER
ENV UID=1000
ENV GID=1000
ENV APP_HOME="/home/cert400/app"

RUN addgroup --gid "$GID" "$GROUPNAME"
RUN adduser --gecos "" --disabled-password --ingroup "$GROUPNAME" --uid "$UID" "$USER"

RUN mkdir "$APP_HOME"

ADD . "$APP_HOME"

WORKDIR "$APP_HOME"

RUN go build -o main .

USER cert400

CMD ["./main", "server"]
