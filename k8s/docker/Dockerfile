FROM golang

WORKDIR /app

COPY authsvc /app/authsvc
COPY apisvc /app/apisvc

EXPOSE 9000
EXPOSE 9001

CMD ["/app/authsvc", "/app/apisvc"]