FROM alpine:3.10

ADD stack_config.yml /bin/stack_config.yml
ADD bin/stackway /bin/stackway

WORKDIR /bin
ENTRYPOINT [ "stackway" ]
CMD [ "--config", "stack_config.yml" ]
