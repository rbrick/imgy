FROM golang:1.8.3-alpine

WORKDIR "/go/src/github.com/rbrick/imgy"

RUN apk add --no-cache curl sqlite git gcc g++

# Install Glide
RUN curl https://glide.sh/get | sh

# Get our dependencies
COPY glide.yaml /go/src/github.com/rbrick/imgy/glide.yaml
COPY glide.lock /go/src/github.com/rbrick/imgy/glide.lock
RUN glide install


# Copy the environment variables
COPY imgy.env /go/src/github.com/rbrick/imgy/imgy.env
COPY setup.sh /go/src/github.com/rbrick/imgy/setup.sh

RUN sh setup.sh

# Copy everything else
COPY . /go/src/github.com/rbrick/imgy/

RUN go build .

# CMD ["./imgy"]

