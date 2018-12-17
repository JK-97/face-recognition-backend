FROM ubuntu:16.04

# bug in ubuntu 16.04 timezone 
# See: https://serverfault.com/questions/683605/docker-container-time-timezone-will-not-reflect-changes/683651
ENV TZ Asia/Shanghai
RUN echo $TZ > /etc/timezone && apt-get update && apt-get install -y tzdata && rm /etc/localtime && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime &&  dpkg-reconfigure -f noninteractive tzdata && apt-get clean

WORKDIR /backend
COPY ./bin/face-recognition-backend .

ENTRYPOINT [ "/backend/face-recognition-backend" ]
# CMD [ "sleep 360000" ]

# # Build Stage
# FROM images.jiangxingai.com:5000/face-recognition-backend:1.11 AS build-stage

# LABEL app="build-face-recognition-backend"
# LABEL REPO="https://gitlab.jiangxingai.com/luyor/face-recognition-backend"

# ENV PROJPATH=/go/src/gitlab.jiangxingai.com/luyor/face-recognition-backend

# # Because of https://github.com/docker/docker/issues/14914
# ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# ADD . /go/src/gitlab.jiangxingai.com/luyor/face-recognition-backend
# WORKDIR /go/src/gitlab.jiangxingai.com/luyor/face-recognition-backend

# RUN make build-alpine

# # Final Stage
# FROM images.jiangxingai.com:5000/face-recognition-backend

# ARG GIT_COMMIT
# ARG VERSION
# LABEL REPO="https://gitlab.jiangxingai.com/luyor/face-recognition-backend"
# LABEL GIT_COMMIT=$GIT_COMMIT
# LABEL VERSION=$VERSION

# # Because of https://github.com/docker/docker/issues/14914
# ENV PATH=$PATH:/opt/face-recognition-backend/bin

# WORKDIR /opt/face-recognition-backend/bin

# COPY --from=build-stage /go/src/gitlab.jiangxingai.com/luyor/face-recognition-backend/bin/face-recognition-backend /opt/face-recognition-backend/bin/
# RUN chmod +x /opt/face-recognition-backend/bin/face-recognition-backend

# # Create appuser
# RUN adduser -D -g '' face-recognition-backend
# USER face-recognition-backend

# ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# CMD ["/opt/face-recognition-backend/bin/face-recognition-backend"]
