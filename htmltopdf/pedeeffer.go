package htmltopdf

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type Pedeeffer struct {
	s3 *s3.S3
}

func New() *Pedeeffer {
	p := &Pedeeffer{}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	}))
	p.s3 = s3.New(sess)

	return p
}

func (p *Pedeeffer) ToBase64(ctx *gin.Context) {
	var request Base64Request
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if request.Url == "" {
		ctx.JSON(400, gin.H{"error": "url is required"})
		return
	}

	data, err := p.getDataBytes(ctx, request.Url)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	buf := bytes.NewBuffer(nil)
	_, err = base64.NewEncoder(base64.StdEncoding, buf).Write(data)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": buf.String()})
}

func (p *Pedeeffer) ToS3(ctx *gin.Context) {
	var request S3Request
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if request.Url == "" {
		ctx.JSON(400, gin.H{"error": "url is required"})
		return
	}

	data, err := p.getDataBytes(ctx, request.Url)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = p.uploadToS3(request.S3Bucket, request.FileName, data)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"success": true})
}

func (p *Pedeeffer) getDataBytes(ctx *gin.Context, page string) ([]byte, error) {
	cdpctx, cancel := chromedp.NewContext(
		ctx,
	)
	defer cancel()
	resultCh := make(chan []byte, 1)
	errorCh := make(chan error, 1)
	err := chromedp.Run(cdpctx,
		chromeTask(cdpctx, page, errorCh, resultCh),
	)

	if err != nil {
		return nil, err
	}

	select {
	case err := <-errorCh:
		return nil, err
	default:
	}

	result := <-resultCh

	return result, nil
}

func (p *Pedeeffer) uploadToS3(bucket, file string, data []byte) error {
	_, err := p.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}

	return nil
}

func chromeTask(ctx context.Context, url string, errorChan chan error, resultChan chan []byte) chromedp.Tasks {
	chromedp.ListenTarget(ctx, func(event interface{}) {
		switch responseReceivedEvent := event.(type) {
		case *network.EventResponseReceived:
			response := responseReceivedEvent.Response
			if strings.HasPrefix(response.URL, url) && response.Status > 299 {
				errorChan <- fmt.Errorf(fmt.Sprintf("error loading page %s: HTTP Code: %d, Extra: %s", url, response.Status, response.StatusText))
			}
		}
	})

	return chromedp.Tasks{
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.Sleep(500 * time.Millisecond),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			do, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}

			resultChan <- do

			return err
		}),
	}
}
