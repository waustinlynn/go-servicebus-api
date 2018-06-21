FROM golang:latest 
RUN mkdir /app 
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/context
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/mendsley/gojwk
RUN go get github.com/waustinlynn/servicebus

ADD . /app/ 
WORKDIR /app 
RUN go build -o main . 
CMD ["/app/main"]