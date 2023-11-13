# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.18.14 as builder

ENV GOPROXY http://proxy.golang.org

RUN mkdir -p /src/projectcollection
WORKDIR /src/projectcollection

# this will cache the npm install step, unless package.json changes
ADD package.json .
RUN npm install --no-progress
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

ADD . .
RUN buffalo build --static -o /bin/app

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/app .

# copy the .ssh directory to the container
COPY .ssh /root/.ssh

# Uncomment to run the binary in "production" mode:
# ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0
ENV GODEBUG=netdns=go

EXPOSE 3000
# setup the cron job with flock to prevent multiple instances of the same job
RUN echo "* * * * * flock -n /tmp/cron.lock /bin/app task processor:application" | crontab -
# start the cron daemon in the background, execute the migration and seed the database, then start the app
CMD crond -f -d 8 & /bin/app migrate; /bin/app task db:seed; /bin/app

# Uncomment to run the migrations before running the binary:
# CMD /bin/app migrate; /bin/app
# CMD exec /bin/app
