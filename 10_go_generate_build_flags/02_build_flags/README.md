## Build flags

Include files with flags by adding `-tags XXX` to build command.
For example,
`go build -o app && ./app` OR `go run .` will output `[A, B, C]`
`go build -o app -tags extra && ./app` OR `go run -tags extra .` will output `[A, B, C, D, E, F]`