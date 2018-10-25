FROM scratch

WORKDIR /backend
COPY ./bin/tf-pose-backend /backend
COPY ./web /backend/web

ENTRYPOINT [ "/backend/tf-pose-backend" ]
CMD [ "serve" ]

# # Build Stage
# FROM images.jiangxingai.com:5000/tf-pose-backend:1.11 AS build-stage

# LABEL app="build-tf-pose-backend"
# LABEL REPO="https://gitlab.jiangxingai.com/luyor/tf-pose-backend"

# ENV PROJPATH=/go/src/gitlab.jiangxingai.com/luyor/tf-pose-backend

# # Because of https://github.com/docker/docker/issues/14914
# ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# ADD . /go/src/gitlab.jiangxingai.com/luyor/tf-pose-backend
# WORKDIR /go/src/gitlab.jiangxingai.com/luyor/tf-pose-backend

# RUN make build-alpine

# # Final Stage
# FROM images.jiangxingai.com:5000/tf-pose-backend

# ARG GIT_COMMIT
# ARG VERSION
# LABEL REPO="https://gitlab.jiangxingai.com/luyor/tf-pose-backend"
# LABEL GIT_COMMIT=$GIT_COMMIT
# LABEL VERSION=$VERSION

# # Because of https://github.com/docker/docker/issues/14914
# ENV PATH=$PATH:/opt/tf-pose-backend/bin

# WORKDIR /opt/tf-pose-backend/bin

# COPY --from=build-stage /go/src/gitlab.jiangxingai.com/luyor/tf-pose-backend/bin/tf-pose-backend /opt/tf-pose-backend/bin/
# RUN chmod +x /opt/tf-pose-backend/bin/tf-pose-backend

# # Create appuser
# RUN adduser -D -g '' tf-pose-backend
# USER tf-pose-backend

# ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# CMD ["/opt/tf-pose-backend/bin/tf-pose-backend"]
