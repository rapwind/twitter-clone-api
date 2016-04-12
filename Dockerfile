FROM golang:1.4.2

MAINTAINER POPPO <hayami_rikuo@cyberagent.co.jp>

# Set Timezone
RUN echo "Asia/Tokyo" > /etc/timezone
RUN dpkg-reconfigure -f noninteractive tzdata

# Set GOPATH/GOROOT environment variables
RUN mkdir -p /go
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

# Set up app
WORKDIR /go/src/github.com/techcampman/twitter-d-server
COPY . /go/src/github.com/techcampman/twitter-d-server/
RUN go get github.com/tools/godep
RUN go get ./...
RUN go build -race -o /go/bin/poppo-api .

# Removed unnecessary packages
RUN apt-get autoremove -y

# Clear package repository cache
RUN apt-get clean all

EXPOSE 3000

CMD ["/go/bin/poppo-api"]
