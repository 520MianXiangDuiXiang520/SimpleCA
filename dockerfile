FROM golang:latest

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct
WORKDIR /SimpleCA
COPY . /SimpleCA
RUN chmod 777 ./run.sh
RUN chmod +x ./bin/simpleCA
RUN make build
EXPOSE 8080
ENTRYPOINT ["./run.sh"] >> simpleCA.log
