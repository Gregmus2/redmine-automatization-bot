FROM golang

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT /entrypoint.sh