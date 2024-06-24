# Go Auth Demo

This is a demo on how to implement authentication in a Go application. It covers basic username and password authentication, jwts, as well as OAuth2 authentication using GitHub.

## Getting Started

To get started, clone the repository and set up your environment variables:

```bash
git clone https://github.com/Cprime50/go-auth.git
cd go-auth
go mod download
```
```bash
export ACCESS_SECRET_KEY='your_access_secret_key_here'
export TOKEN_TTL=5m
export GITHUB_CLIENT_ID='your_github_client_id'
export GITHUB_CLIENT_SECRET='your_github_client_secret'
```
```bash
go run main.go
```