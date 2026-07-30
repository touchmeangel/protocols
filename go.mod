module github.com/touchmeangel/protocols

go 1.25.0

require (
	github.com/touchmeangel/twitter_api v0.0.11
	github.com/valyala/fasthttp v1.73.0
	golang.org/x/net v0.57.0
)

require (
	github.com/andybalholm/brotli v1.2.2 // indirect
	github.com/klauspost/compress v1.19.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/text v0.40.0 // indirect
)

replace golang.org/x/net => github.com/golang/net v0.57.0

replace golang.org/x/text => github.com/golang/text v0.40.0

replace golang.org/x/crypto => github.com/golang/crypto v0.48.0

replace golang.org/x/sys => github.com/golang/sys v0.47.0
