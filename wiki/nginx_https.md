# 在Nginx下以HTTPS方式运行 #

    server {
        listen          		80;
        server_name     		*.ts.com;
        client_max_body_size  	10m;
    	rewrite ^(.*)$ https://$host$1 permanent; 
    }
   
    server {
        listen          		443 ssl;
        server_name     		master.ts.com super.ts.com;
		ssl_certificate      	/home/usr/go/src/go2o/ssl/ssl_certificate.cer;
        ssl_certificate_key  	/home/usr/go/src/go2o/ssl/ssl_certificate.key;
        client_max_body_size  	10m;
	    location / {
            proxy_pass   http://localhost:14281;
            proxy_set_header Host $host;
        }
    }

    server {
        listen          		443 ssl;
        server_name     		static.ts.com;
  		ssl_certificate      	/home/usr/go/src/go2o/ssl/ssl_certificate.cer;
    	ssl_certificate_key  	/home/usr/go/src/go2o/ssl/ssl_certificate.key;
        root    	    		/home/usr/go/src/go2o/static/;
	 	location / {
            expires 1h;    
        }
	    location ~* \.(eot|ttf|woff|woff2|svg)$ {
  			add_header Access-Control-Allow-Origin *;
			expires 10d;
      	}
    }

    server {
        listen          		443 ssl;
        server_name     		img.ts.com;
        ssl_certificate      	/home/usr/go/src/go2o/ssl/ssl_certificate.cer;
        ssl_certificate_key  	/home/usr/go/src/go2o/ssl/ssl_certificate.key;
        root            		/home/usr/go/src/go2o/uploads;
	    location / {
      	  	expires 1d;
        }  
	  	location ~* \.(eot|ttf|woff|woff2|svg)$ {
      		add_header Access-Control-Allow-Origin *;
      	}
    }

  
    server{
		listen    443 ssl;
        server_name          	*.ts.com localhost;
        ssl_certificate      	/home/usr/go/src/go2o/ssl/ssl_certificate.cer;
		ssl_certificate_key  	/home/usr/go/src/go2o/ssl/ssl_certificate.key;
        client_max_body_size  	10m;	 
		location ~*\.txt$ {
            root /home/usr/go/src/go2o/;
        }
		location / {
          proxy_pass   http://localhost:14190;
          proxy_set_header Host $host;
        }
   }
