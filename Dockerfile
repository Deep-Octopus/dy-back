FROM busybox
MAINTAINER OCTO
WORKDIR /root
ADD ./dy ./dy
ADD ./config ./config
ADD ./assets ./assets
EXPOSE 8888
ENTRYPOINT  ["./dy"]