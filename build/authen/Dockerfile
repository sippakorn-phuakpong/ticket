FROM cosmtrek/air AS development

WORKDIR /app

COPY cmd/authen ./
COPY api/authen api/authen
COPY pkg pkg
COPY go.mod go.sum ./
COPY config config
COPY certs certs
COPY build/authen ./

RUN go mod tidy
RUN go build -v -o ./dist/authen

CMD ["air", "-c", ".air.toml"]

FROM debian:bookworm AS production

ARG timezone=Asia/Bangkok
ARG hotreload=false
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV TZ $timezone

WORKDIR /app

RUN apt-get update && apt-get install

COPY --from=development /app/dist/authen .
COPY --from=development  /app/config config
COPY --from=development  /app/certs certs

RUN chmod +x authen

CMD ["./authen"]