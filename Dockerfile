FROM golang:latest
EXPOSE 8080

ENV GO111MODULE=on
RUN mkdir exam_fpt
COPY . /exam_fpt/
WORKDIR /exam_fpt
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build


ENTRYPOINT ["./exam_fpt"]

