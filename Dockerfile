FROM golang:1.11

LABEL maintainer="Harry Brown <harrybrown98@gmail.com"

RUN mkdir src/harrybrown.com
COPY . src/harrybrown.com

# RUN env
RUN ls
RUN ls src/harrybrown.com
# RUN cd harrybrown.com
RUN ls harrybrown.com

RUN go get -d -v harrybrown.com/...
# RUN go build -o harrybrown.com/harrybrown.com main.go
# RUN go get ./...
RUN go install harrybrown.com

EXPOSE 8080
CMD ["harrybrown.com"]