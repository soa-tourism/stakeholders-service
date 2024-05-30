FROM golang:latest
WORKDIR /go/stakeholders/stakeholders-service
RUN go mod init stakeholders-service
COPY stakeholders .
RUN go build -o main .
EXPOSE 81
CMD ["./main"]