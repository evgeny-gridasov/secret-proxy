#### Overview
Secret Proxy is a small web server that starts connection forwarding on-demand. Could be useful to temporary forward one-time connection from a public server into a private network.

A single incoming connection is allowed within a 5-second time interval, after which, the listener will stop. That should be more than enough to click the button in a web browser and initiate client connection.

It is recommended to put Secret Proxy behind an NGINX proxy, protected by TLS and http authentication.

#### Building

To build, make sure you have Golang installed, download secret-proxy and type:

```
go build

```

There are no external dependencies, so running command above should create the `secret-proxy` executable:
```
$ ./secret-proxy 
secret-proxy is a small web server that starts connection forwarding on-demand

Usage:

  -from string
    	From host:port
  -to string
    	To host:port
  -web string
    	Web server host:port

``` 

#### Credits
Written by Evgeny Gridasov ([@evgeny-gridasov](http://github.com/evgeny-gridasov))