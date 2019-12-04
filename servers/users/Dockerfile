FROM alpine
RUN apk add --no-cache ca-certificates 
COPY users /users
ENTRYPOINT ["/users"]