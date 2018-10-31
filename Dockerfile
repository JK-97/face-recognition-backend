FROM busybox AS build-env
RUN touch /empty

FROM scratch

WORKDIR /backend
COPY ./bin/tf-fence-backend .
COPY ./web ./web
COPY --from=build-env /empty /img

ENTRYPOINT [ "/backend/tf-fence-backend" ]
CMD [ "serve" ]

# # Build Stage
# FROM images.jiangxingai.com:5000/tf-fence-backend:1.11 AS build-stage

# LABEL app="build-tf-fence-backend"
# LABEL REPO="https://gitlab.jiangxingai.com/luyor/tf-fence-backend"

# ENV PROJPATH=/go/src/gitlab.jiangxingai.com/luyor/tf-fence-backend

# # Because of https://github.com/docker/docker/issues/14914
# ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# ADD . /go/src/gitlab.jiangxingai.com/luyor/tf-fence-backend
# WORKDIR /go/src/gitlab.jiangxingai.com/luyor/tf-fence-backend

# RUN make build-alpine

# # Final Stage
# FROM images.jiangxingai.com:5000/tf-fence-backend

# ARG GIT_COMMIT
# ARG VERSION
# LABEL REPO="https://gitlab.jiangxingai.com/luyor/tf-fence-backend"
# LABEL GIT_COMMIT=$GIT_COMMIT
# LABEL VERSION=$VERSION

# # Because of https://github.com/docker/docker/issues/14914
# ENV PATH=$PATH:/opt/tf-fence-backend/bin

# WORKDIR /opt/tf-fence-backend/bin

# COPY --from=build-stage /go/src/gitlab.jiangxingai.com/luyor/tf-fence-backend/bin/tf-fence-backend /opt/tf-fence-backend/bin/
# RUN chmod +x /opt/tf-fence-backend/bin/tf-fence-backend

# # Create appuser
# RUN adduser -D -g '' tf-fence-backend
# USER tf-fence-backend

# ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# CMD ["/opt/tf-fence-backend/bin/tf-fence-backend"]
