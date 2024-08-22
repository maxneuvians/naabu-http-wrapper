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

FROM mcr.microsoft.com/devcontainers/go:1-1.21-bullseye

# Install packages
RUN apt-get update
RUN apt-get install -y libpcap-dev ca-certificates

COPY --from=build /naabu-http-wrapper /naabu-http-wrapper

ENTRYPOINT ["/naabu-http-wrapper"]
