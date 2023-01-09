FROM golang:1.17

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/codefresh-contrib/go-sample-app

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...


# This container exposes port 8080 to the outside world
EXPOSE 8080

RUN go build -o ./out/go-sample-app .

# Run the executable
CMD ["./out/go-sample-app"]