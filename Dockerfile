# build stage
FROM golang:1.19 as builder

ARG PROGRAM_VER=dev-docker
ENV CGO_ENABLED=0

RUN apt-get -qq update && \
	apt-get install -yqq upx

COPY . /build
WORKDIR /build

RUN go build -ldflags "-X main.programVer=${PROGRAM_VER}" -o /build/app
RUN strip /build/app
RUN upx -q -9 /build/app

# ---
FROM scratch

ENV TZ=Asia/Seoul

ENV MQTT_HOST="secret"
ENV MQTT_USERNAME="secret"
ENV MQTT_PASSWORD="secret"
ENV TELEGRAM_APITOKEN="secret"
ENV TELEGRAM_ROOM_ID="secret"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/app .

# Diag http port
EXPOSE 8080

ENTRYPOINT ["./app"]