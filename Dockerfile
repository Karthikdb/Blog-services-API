FROM golang:1.11

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/lib
RUN git clone https://github.com/lib/pq.git
WORKDIR $GOPATH/src/github.com/rs
RUN git clone https://github.com/rs/cors.git
# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
WORKDIR $GOPATH/src
COPY . .
RUN go get github.com/lib/pq
RUN go get github.com/rs/cors
# Download all the dependencies
RUN go build
# Install the package


# This container exposes port 8080 to the outside world


# Run the executable
CMD ["./src"]
