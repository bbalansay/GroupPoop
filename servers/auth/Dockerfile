FROM alpine
RUN apk add --no-cache ca-certificates 
COPY auth /auth
ENTRYPOINT ["/auth"]