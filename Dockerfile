FROM golang:1.17.2-alpine3.14

ENV DIR_HOME=/root
ENV DIR_DATA /data

WORKDIR ${DIR_DATA}

# copy the repo
COPY . reminder

# install the command
RUN cd reminder \
    && go install cmd/reminder/main.go

# rename the command
RUN cp ${GOPATH}/bin/main ${GOPATH}/bin/reminder

WORKDIR ${DIR_HOME}

CMD [ \
        "/bin/sh", "-c", \
        " \
        reminder \
        # while true; do echo \"Hit CTRL+C\"; sleep 1; done \
        " \
    ]
