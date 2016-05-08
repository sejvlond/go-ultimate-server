FROM scratch
MAINTAINER Ondřej Šejvl

VOLUME [ "/www/ultimate-server/conf" ]

WORKDIR /www/ultimate-server/

# default configuration
COPY conf/* ./conf/
COPY go-ultimate-server ./bin/server
COPY templates ./templates/

EXPOSE 9876

ENTRYPOINT [ "/www/ultimate-server/bin/server" ]
CMD [ "-c", "/www/ultimate-server/conf/conf.yaml" ]

