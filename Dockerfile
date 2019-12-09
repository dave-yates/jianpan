#our base image
FROM golang:1.8.3 as builder 

#set the workdir
WORKDIR src/github.com/dave-yates/jianpan

#copy the code to the new container
COPY / .

#compile the go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/src/github.com/dave-yates/jianpan .

#port number to be exposed
#EXPOSE 8080

#run the server
CMD ["./main"]


#instructions
#docker build -t keyboard .
#docker run -ti -p 127.0.0.1:8080:8080 keyboard
#go to http://localhost:8080/keyboard for the keyboard
#go to http://localhost:8080/help for a few hints
