# Stage 1
FROM golang:1.18-alpine as builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app/go

COPY . .

# Install required packages
RUN go mod download

# Build the program
RUN GOOS=linux GOARCH=amd64 go build ./main.go

# Stage 2
FROM golang:1.18-alpine

COPY --from=builder /app/go/main /home/main

# Set the start command
CMD ["/home/main"]
