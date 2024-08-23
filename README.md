# Naabu HTTP wrapper
A simplistic HTTP wrapper around Naabu to allow it to be invoked through HTTP. You can pass it the following GET params

```
host = Host to scan ex. scanme.sh (required)
ports = Ports to scan ex. 80,443 or 1-1000 (optional)
timeouts = Timeout in milliseconds (default: 1000)
```

The webserver will run Naabu and return a JSON response that looks like this:

```
{"Host":"scanme.sh","Ports":[80,443]}
```

## Usage

The following commands will:

1. Build a docker image with the webserver
2. Run the docker container locally with webserver running on port 8080
3. Query the webserver using curl

```
make build
make run
curl 'http://127.0.0.1:8080/?host=scanme.sh&ports=80,443'
```

## Building binary locally

Build your local binary like any other Golang binary, however, Naabu requires that you have `libpcap-devel` installed on your system.

## License
MIT 2024