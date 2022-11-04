# DNS Server with TLS and Plugins

The idea is to offer a single executable DNS server that can be configured to provide multiple transport mechanisms and can be extended with plugins. It supports the following transports:

- UDP port 53
- TCP port 53
- DNS over TLS (DoT) on port 853
- DNS over HTTPS (DoH) on port 443

Each listener can be configured separately and can be enabled or disabled.

## Plugins

The server can be configured to use plugins. The following plugins are available:

- `Proxy` - resolves queries from an upstream DNS server
- `Stats` - gathers statistics about the queries (number of lookups per domain, log queries)
- `Logger` - log queries and responses
- `Blocker` - blocks queries using a block list

## Getting started

As the idea is to have a single file DNS server, the server is written in Go and can be compiled to a single binary. That's why we do not have an external config file. Instead, the server is configured directly in the source code.

You can adjust the configuration under `config/dev.go`. The configuration is pretty self-explanatory. You can enable or disable the listeners and plugins. You can also configure the plugins.

To run the server, you can run the following command:

```
make run
```

The server should now be running on your machine. You can test it by running the following command:

```
$ dig @127.0.0.1 example.com

; <<>> DiG 9.10.6 <<>> @127.0.0.1 example.com
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 42396
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;example.com.                   IN      A

;; ANSWER SECTION:
example.com.            85931   IN      A       93.184.216.34

;; Query time: 84 msec
;; SERVER: 127.0.0.1#53(127.0.0.1)
;; WHEN: Thu Nov 03 22:39:49 CET 2022
;; MSG SIZE  rcvd: 56
```

If you want to test the DNS over TLS (DoT) or DNS over HTTPS (DoH) listener, I highly recommend using [doggo](https://github.com/mr-karan/doggo) which is a super simple and versatile DNS client. You can query the DoT listener with the following command:

```
$ doggo @tls://localhost example.com

NAME            TYPE    CLASS   TTL     ADDRESS         NAMESERVER
example.com.    A       IN      77754s  93.184.216.34   localhost:853
```

or for DoH:

```
$ doggo @https://localhost/dns-query example.com

NAME            TYPE    CLASS   TTL     ADDRESS         NAMESERVER
example.com.    A       IN      80180s  93.184.216.34   https://localhost/dns-query
```

## Build & Deployment

To build the server for your current platform, you can run the following command:

```
make build
```

To build it for Linux, you can run the following command (make sure you have the required compile chain installed):

```
make build-linux-amd64
```

This gives you a single binary called `dns` that you can run on your server.
