#### Overview
Secret Proxy is a small web server that starts connection forwarding on-demand. Could be useful to temporary forward one-time connection from a public server into a private network.

Only one incoming connection will be allowed within a 5-second time interval, after which, the listener will stop. That should be more than enough to click the button in a web browser to start connection forwarding and initiate client connection.

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

#### Usage

Run secret-proxy either in a screen/tmux session or as a systemd/init.d service. For instance, if you want to forward a connection to your home ssh server:
```
$ ./secret-proxy -from :9000 -to 192.168.1.100:22 -web localhost:8080
2020/08/30 13:00:04 Web server listening on localhost:8080
```

Set up nginx proxy, configure credentials in `htpasswd` file and use `certbot` to set up lets encrypt certs:
```
server {
        listen 443 ssl;
        listen [::]:443 ssl;

        root /usr/share/nginx/html;
        index index.html index.htm;

        server_name secret.example.com;

        location / {
           expires -1;
           auth_basic "Secret Area";
           auth_basic_user_file htpasswd;
           proxy_pass http://127.0.0.1:8080/;
        }

        error_page 404 /404.html;

        error_page 500 502 503 504 /50x.html;
        location = /50x.html {
                root /usr/share/nginx/html;
        }
        ssl_certificate /etc/letsencrypt/live/secret.example.com/fullchain.pem; # managed by Certbot
        ssl_certificate_key /etc/letsencrypt/live/secret.example.com/privkey.pem; # managed by Certbot
}

```
Next, navigate to https://secret.example.com, authenticate with credentials configured in `htpasswd` and click "Start Forwarding". Within 5 seconds, ssh onto your public server to port 9000:
```
ssh -p 9000 secret.example.com

```
The ssh connection will be forwarded to 192.168.1.100:22 and the listener on port 9000 will be stopped. Same works great with RDP forwarded connections, too.

#### Credits
Written by Evgeny Gridasov ([@evgeny-gridasov](http://github.com/evgeny-gridasov))