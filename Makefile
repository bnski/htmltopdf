default: clean html2pdf_darwin_arm64

clean:
	go clean
	rm -f html2pdf_darwin_arm64
	rm -f html2pdf_linux_arm64
	rm -f html2pdf_x86_64

html2pdf_darwin_arm64:
	CGO_ENABLED=0 go build -o html2pdf_darwin_arm64 -ldflags "-s -w"

html2pdf_linux_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o html2pdf_linux_arm64 -ldflags "-s -w"

html2pdf_x86_64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o html2pdf_x86_64 -ldflags "-s -w"
