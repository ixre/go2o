module github.com/jsix/go2o

require (
	git.apache.org/thrift.git v0.11.0
	github.com/afocus/captcha v0.0.0-20170421134744-6d694b359d1a
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jsix/alidayu v0.0.0-20160128071321-7eadea36c79c
	github.com/jsix/goex v1.0.5
	github.com/jsix/gof v1.0.0
	github.com/labstack/echo v3.2.1+incompatible
	github.com/labstack/gommon v0.2.1 // indirect
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.3 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/robfig/cron v0.0.0-20180505203441-b41be1df6967
	github.com/stretchr/testify v1.2.2 // indirect
	github.com/valyala/bytebufferpool v0.0.0-20160817181652-e746df99fe4a // indirect
	github.com/valyala/fasttemplate v0.0.0-20170224212429-dcecefd839c4 // indirect
	golang.org/x/text v0.3.0
	gopkg.in/square/go-jose.v1 v1.1.2
)

replace (
	git.apache.org/thrift.git v0.11.0 => github.com/TriangleGo/thrift v0.11.0
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
)
