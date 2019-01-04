FROM golang:latest
RUN mkdir /app
ADD *.go /app/
WORKDIR /app
RUN go build -o dwarves .
CMD ["/app/dwarves"]
