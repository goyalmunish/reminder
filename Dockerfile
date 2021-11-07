FROM golang:1.17.2-alpine3.14

ENV DIR_HOME=/app
ENV DIR_DATA /data

# install required packages
RUN apk add git
RUN apk add openssh

# clone the repo
WORKDIR ${DIR_DATA}
RUN git clone https://github.com/goyalmunish/reminder.git

# install the command
RUN cd reminder \
    && go install cmd/reminder/main.go

# rename the command
RUN cp ${GOPATH}/bin/main ${GOPATH}/bin/reminder

WORKDIR ${DIR_HOME}

CMD [ \
        "/bin/sh", "-c", \
        " \
        while true; do echo \"Hit CTRL+C\"; sleep 1; done \
        " \
    ]
