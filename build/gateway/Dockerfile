FROM cosmtrek/air AS development

WORKDIR /app

COPY cmd/gateway ./

COPY pkg pkg
COPY go.mod go.sum ./
COPY config config
COPY certs certs
COPY build/gateway ./

RUN go mod tidy
RUN go build -v -o ./dist/gateway

CMD ["air", "-c", ".air.toml"]

FROM debian:bookworm  AS production

ARG timezone=Asia/Bangkok
ARG hotreload=false
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV TZ $timezone

WORKDIR /app

RUN apt-get update && apt-get install

COPY --from=development /app/dist/gateway .
COPY --from=development  /app/config config

RUN chmod +x gateway

CMD ["./gateway"]