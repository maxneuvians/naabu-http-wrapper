FROM mcr.microsoft.com/devcontainers/go:1-1.21-bullseye AS build

# Install packages
RUN apt-get update
RUN apt-get install -y libpcap-dev

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -o /naabu-http-wrapper

FROM bitnami/minideb:latest

ENV USER=app
ENV UID=10001

# Create non-root user
RUN useradd -m -s /bin/bash -u $UID $USER

# Install packages
RUN apt-get update
RUN apt-get install -y libpcap-dev ca-certificates

COPY --from=build --chown=${USER}:${USER} /naabu-http-wrapper /home/app/naabu-http-wrapper

USER ${USER}:${USER}

ENTRYPOINT ["/home/app/naabu-http-wrapper"]
