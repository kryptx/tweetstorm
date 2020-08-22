FROM golang:latest

WORKDIR /src
COPY . .
RUN go build -o /out/tweetstorm .

EXPOSE 3000
CMD /out/tweetstorm
