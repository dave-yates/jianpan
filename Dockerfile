#our base image
FROM golang:latest

#set the workdir
WORKDIR src/github.com/dave-yates/jianpan

#copy the code to the new container
COPY / .

RUN go get -v -u go.mongodb.org/mongo-driver/mongo

#port number to be exposed
EXPOSE 8080

#run the server
CMD ["go", "run",  "main.go"]

