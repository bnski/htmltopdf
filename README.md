### Notes to Self:

**Install Go**
```
brew install go
```

**Packages/Modules/Cheats**
```
go mod tidy
go mod vendor 
go clean -modcache   
lsof -i :PORT
kill -9 PID
```

**Build**
```
make html2pdf_darwin_arm64
make html2pdf_linux_arm64       
make html2pdf_x86_64  
```

**Web Server**
.env
```
PORT=8080
```

```
go run main.go   
```

http://localhost:8080


**The API exposes 2 endpoints:**
1. POST /page-to-base64
2. POST /page-to-s3

example base64 requests:
```json
{
  "url": "https://www.google.com"
}
```

example s3 requests:
```json
{
  "url": "https://www.google.com",
  "bucket": "bucket-name",
  "file_name": "file-name.pdf"
}
```
