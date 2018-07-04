FROM alpine:3.6

RUN apk add --no-cache --virtual .build-deps \
        tzdata \
        && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
        && echo "Asia/Shanghai" > /etc/timezone \
        && apk del .build-deps

ENV TZ "Asia/Shanghai"

COPY ./snowflake /bin/snowflake

EXPOSE 11070

VOLUME ["/etc/snowflake"]

CMD ["snowflake", "--config", "/etc/snowflake/config.yaml"]
