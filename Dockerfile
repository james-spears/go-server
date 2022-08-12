# Stage 1
FROM golang:1.19-alpine as builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /opt/app

COPY . .

# Install required packages
RUN go mod download

# Build the program
RUN GOOS=linux GOARCH=amd64 go build ./main.go

# Stage 2
FROM alpine:3.16

COPY --from=builder /opt/app/main /opt/app/main

# Set the start command
CMD ["/opt/app/main"]
