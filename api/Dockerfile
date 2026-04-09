FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY foresight-api .

ENV DB_HOST=10.255.12.75
ENV DB_USER=root
ENV DB_PASS=Gonzo@7072
ENV DB_NAME=ads
ENV APP_PORT=10000
ENV APP_HOST=localhost
EXPOSE $APP_PORT
EXPOSE 3306

CMD ["/app/foresight-api"]