# tiny-wiki
A wiki system based on mdwiki, and has functions to auto generate wiki directory。

#### Run generate directory
```bash
go run generate.go
```

#### Run in nginx
```bash
server {
	listen 80;
	server_name gavin.com;
    location / {
        root   /data/www/tiny-wiki;
        index  index.html index.htm index.php;
    }
}
```