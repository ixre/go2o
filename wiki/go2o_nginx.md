创建配置文件
```
server {
    listen  80;
    listen  443 ssl;
    server_name api.xxx.com;
    ssl_certificate      /data/cert/ssl_certificate.cer;
    ssl_certificate_key  /data/cert/ssl_certificate.key;
    ssl_protocols        SSLv3 SSLv2 TLSv1 TLSv1.1 TLSv1.2;
    ssl_session_cache    shared:SSL:1m;
    ssl_session_timeout  5m
    client_max_body_size 10m;
    location / {
        proxy_pass   http://localhost:1428;
        proxy_set_header Host $host;
    }
}

server{
    listen                  80;
    listen          		443 ssl;
    server_name          	*.xxx.com localhost;
    client_max_body_size  	10m;
    location / {
      proxy_pass   http://localhost:14190;
      proxy_set_header Host $host;
    }
}

server {
    listen                  80;
    listen          		443 ssl;
    server_name     		master.xxx.com;
    client_max_body_size  	10m;
    location / {
        proxy_pass   http://localhost:14281;
        proxy_set_header Host $host;
    }
}
    
server {
    listen                  80;
    listen          		443 ssl;
    server_name     		static.xxx.com;
    root    	    		/go2o/static/;
    location / {
        expires 1h;
    }
    location ~* \.(eot|ttf|woff|woff2|svg)$ {
        add_header Access-Control-Allow-Origin *;
        expires 10d;
    }
}

server {
    listen                  80;
    listen          		443 ssl;
    server_name     		img.xxx.com;
    root            		/go2o/uploads;
    location / {
        expires 1d;
    }
    location ~* \.(eot|ttf|woff|woff2|svg)$ {
        add_header Access-Control-Allow-Origin *;
    }
}
```
_注:将xxx.com替换为你的域名,包括证书_

可以添加以下配置, 强制重定向到HTTPS
```
server{
    listen          80;
    server_name     *.xxx.com;
    rewrite ^(.*)$ https://$host$1 permanent;
}
```