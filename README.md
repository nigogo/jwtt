# jwtt

**jwtt** (JWT tool) is a simple CLI application for decoding JWT tokens.

> **Disclaimer** This is a messy, incomplete and untested project. Do not use it in production. It will let you inspect the contents of a JWT token, and nothing more.

## Usage

Build the project:

```bash
go build -o jwtt main.go
```

Feed the JWT token to the `jwtt` command (use the `-t` flag to convert timestamps to a readable format):

```bash
jwtt -t eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

The output will be the decoded token:

```plaintext
Header:
{
  "alg": "HS256",
  "typ": "JWT"
}
Payload:
{
  "iat": "1516239022 (2018-01-18T02:30:22+01:00)",
  "name": "John Doe",
  "sub": "1234567890"
}
Signature:
SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```
