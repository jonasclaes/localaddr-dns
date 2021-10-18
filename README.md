# localaddr-dns
localaddr-dns is a private network DNS reflector written in Golang.

# How to use
- Download the repo by running `git clone https://github.com/jonasclaes/localaddr-dns.git`
- Build the binary by running `go build` in the source code directory
- Run the binary by running `./localaddr-dns`

For options you can run the binary using the `-h` option.

# Example usage
`localaddr-dns -port 25 -ttl 86400 -base_domain example.com`

# License
This project is licensed under the Apache License 2.0.