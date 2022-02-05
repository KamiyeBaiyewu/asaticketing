FROM golang:1.14.6 as builder 

ENV GO111MODULE=on

ENV MODULE=github.com/lilkid3/ASA-Ticket/Backend
ENV DATA_DIRECTORY /${MODULE}/
WORKDIR ${DATA_DIRECTORY}
ARG APP_VERSION
ENV APP_VERSION $APP_VERSION
ARG CGO_ENABLED=0
COPY . .
# RUN ls internal/config
# RUN pwd
# RUN echo $GOPATH
RUN go build -ldflags="-X '${MODULE}/internal/config.Version=${APP_VERSION}'"  ${DATA_DIRECTORY}cmd/server
# RUN go tool nm ./server | grep internal/config
# RUN go build cmd/server


FROM alpine:3.10
ENV DATA_DIRECTORY=/github.com/lilkid3/ASA-Ticket/Backend/
# ENV APP_VERSION APP_VERSION
RUN apk add --update --no-cache \
    ca-certificates
COPY --from=builder ${DATA_DIRECTORY}/internal/storage/database/migrations ${DATA_DIRECTORY}/internal/storage/database/migrations
COPY --from=builder ${DATA_DIRECTORY}/config /config

COPY --from=builder ${DATA_DIRECTORY}server /server

ENTRYPOINT [ "/server" ]