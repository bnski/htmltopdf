package htmltopdf

type Base64Request struct {
	Url string `json:"url"`
}

type S3Request struct {
	Url      string `json:"url"`
	S3Bucket string `json:"bucket"`
	FileName string `json:"file_name"`
}
