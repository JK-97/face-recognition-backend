FROM ubuntu:16.04

# bug in ubuntu 16.04 timezone 
# See: https://serverfault.com/questions/683605/docker-container-time-timezone-will-not-reflect-changes/683651
ENV TZ Asia/Shanghai
RUN echo $TZ > /etc/timezone && apt-get update && apt-get install -y tzdata && rm /etc/localtime && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime &&  dpkg-reconfigure -f noninteractive tzdata && apt-get clean

WORKDIR /backend
COPY ./bin/face-recognition-backend .

ENTRYPOINT [ "/backend/face-recognition-backend", "serve"]
