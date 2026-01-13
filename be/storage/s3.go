package storage

import (
	"be/model"
	"be/constant"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

var (
	svc *s3.S3
	ErrS3NotInitialized = fmt.Errorf("S3 service not initialized")
)

func S3Init() error {
    if svc != nil {
        return nil // Already initialized
    }

    if constant.APPCONFIG.S3Endpoint == "" {
        return fmt.Errorf("S3 endpoint not configured")
    }

    // Initialize AWS session using MinIO endpoint
    sess, err := session.NewSession(&aws.Config{
        Region:           aws.String(constant.APPCONFIG.S3Region),
        Endpoint:         aws.String(constant.APPCONFIG.S3Endpoint),
        S3ForcePathStyle: aws.Bool(true),
        Credentials:      credentials.NewStaticCredentials(constant.APPCONFIG.S3ID, constant.APPCONFIG.S3Secret, constant.APPCONFIG.S3Token),
    })
    if err != nil {
        return fmt.Errorf("failed to create session: %v", err)
    }

    svc = s3.New(sess)
    
    // Test connection
    _, err = svc.ListBuckets(&s3.ListBucketsInput{})
    if err != nil {
        svc = nil // Reset service on failure
        return fmt.Errorf("failed to connect to S3: %v", err)
    }

    return nil
}

func checkS3Service() error {
    if svc == nil {
        return ErrS3NotInitialized
    }
    return nil
}

func GeneratePresignedURL(bucket, key string, expiry time.Duration) (string, error) {
    if err := checkS3Service(); err != nil {
        return "", err
    }

    if bucket == "" || key == "" {
        return "", fmt.Errorf("bucket and key cannot be empty")
    }

    req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    })
    
    url, err := req.Presign(expiry)
    if err != nil {
        return "", fmt.Errorf("failed to generate pre-signed URL: %v", err)
    }

    return url, nil
}

func UploadSpeechMarks(bucket string, key string, speechMarksData []model.SpeechMarkData) error {
    if err := checkS3Service(); err != nil {
        return err
    }

    if bucket == "" || key == "" {
        return fmt.Errorf("bucket and key cannot be empty")
    }

    if len(speechMarksData) == 0 {
        return fmt.Errorf("speech marks data is empty")
    }

    encodedData, err := json.Marshal(speechMarksData)
    if err != nil {
        return fmt.Errorf("failed to marshal SpeechMarksData: %v", err)
    }

    file := bytes.NewReader(encodedData)
    input := &s3.PutObjectInput{
        Bucket:      aws.String(bucket),
        Key:         aws.String(key),
        Body:        file,
        ContentType: aws.String("application/json"),
    }

    _, err = svc.PutObject(input)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case s3.ErrCodeNoSuchBucket:
                return fmt.Errorf("bucket does not exist: %v", aerr)
            default:
                return fmt.Errorf("failed to upload SpeechMarksData: %v", aerr)
            }
        }
        return fmt.Errorf("failed to upload SpeechMarksData: %v", err)
    }

    return nil
}

func UploadText(bucket string, key string, text string) error {
    if err := checkS3Service(); err != nil {
        return err
    }

    if bucket == "" || key == "" {
        return fmt.Errorf("bucket and key cannot be empty")
    }

    if text == "" {
        return fmt.Errorf("text content is empty")
    }

    file := bytes.NewReader([]byte(text))
    input := &s3.PutObjectInput{
        Bucket:      aws.String(bucket),
        Key:         aws.String(key),
        Body:        file,
        ContentType: aws.String("text/plain"),
    }

    _, err := svc.PutObject(input)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case s3.ErrCodeNoSuchBucket:
                return fmt.Errorf("bucket does not exist: %v", aerr)
            default:
                return fmt.Errorf("failed to upload text: %v", aerr)
            }
        }
        return fmt.Errorf("failed to upload text: %v", err)
    }

    return nil
}

func UploadBase64MP3ToS3(bucket, key, base64Data string) error {
    if err := checkS3Service(); err != nil {
        return err
    }

    if bucket == "" || key == "" {
        return fmt.Errorf("bucket and key cannot be empty")
    }

    if base64Data == "" {
        return fmt.Errorf("base64 data is empty")
    }

    // Remove the prefix (if present)
    const prefix = "data:audio/mp3;base64,"
    if strings.HasPrefix(base64Data, prefix) {
        base64Data = base64Data[len(prefix):]
    }

    decodedData, err := base64.StdEncoding.DecodeString(base64Data)
    if err != nil {
        return fmt.Errorf("failed to decode base64 data: %v", err)
    }

    if len(decodedData) == 0 {
        return fmt.Errorf("decoded MP3 data is empty")
    }

    dataReader := bytes.NewReader(decodedData)
    input := &s3.PutObjectInput{
        Bucket:      aws.String(bucket),
        Key:         aws.String(key),
        Body:        dataReader,
        ContentType: aws.String("audio/mp3"),
    }

    _, err = svc.PutObject(input)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case s3.ErrCodeNoSuchBucket:
                return fmt.Errorf("bucket does not exist: %v", aerr)
            default:
                return fmt.Errorf("failed to upload MP3: %v", aerr)
            }
        }
        return fmt.Errorf("failed to upload MP3: %v", err)
    }

    return nil
}

func Read(bucket, key, downloadPath string) error {
    if err := checkS3Service(); err != nil {
        return err
    }

    if bucket == "" || key == "" || downloadPath == "" {
        return fmt.Errorf("bucket, key, and download path cannot be empty")
    }

    result, err := svc.GetObject(&s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case s3.ErrCodeNoSuchKey:
                return fmt.Errorf("object not found: %v", aerr)
            case s3.ErrCodeNoSuchBucket:
                return fmt.Errorf("bucket does not exist: %v", aerr)
            default:
                return fmt.Errorf("failed to download file: %v", aerr)
            }
        }
        return fmt.Errorf("failed to download file: %v", err)
    }
    defer result.Body.Close()

    outFile, err := os.Create(downloadPath)
    if err != nil {
        return fmt.Errorf("failed to create download file: %v", err)
    }
    defer outFile.Close()

    _, err = outFile.ReadFrom(result.Body)
    if err != nil {
        return fmt.Errorf("failed to copy data to file: %v", err)
    }

    return nil
}

func Delete(bucket, key string) error {
    if err := checkS3Service(); err != nil {
        return err
    }

    if bucket == "" || key == "" {
        return fmt.Errorf("bucket and key cannot be empty")
    }

    _, err := svc.DeleteObject(&s3.DeleteObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case s3.ErrCodeNoSuchKey:
                return fmt.Errorf("object not found: %v", aerr)
            case s3.ErrCodeNoSuchBucket:
                return fmt.Errorf("bucket does not exist: %v", aerr)
            default:
                return fmt.Errorf("failed to delete object: %v", aerr)
            }
        }
        return fmt.Errorf("failed to delete object: %v", err)
    }

    return nil
}

func AutoDelete(ctx context.Context, bucket string) {
    if err := checkS3Service(); err != nil {
        log.Error("AutoDelete failed to start: ", err)
        return
    }

    if bucket == "" {
        log.Error("AutoDelete failed: bucket name is empty")
        return
    }

    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
                Bucket: aws.String(bucket),
            })
            if err != nil {
                log.Error("Error listing objects: ", err)
                continue
            }

            for _, item := range resp.Contents {
                if time.Since(*item.LastModified) > 30*24*time.Hour {
                    if err := Delete(bucket, *item.Key); err != nil {
                        log.Error("Failed to delete ", *item.Key, ": ", err)
                    }
                }
            }
        }
    }
}

func StoreData(requestID, audioData, LLMResp, cleanedMarkdown string, speechMarksData []model.SpeechMarkData) error {
    if err := checkS3Service(); err != nil {
        return err
    }

    if requestID == "" {
        return fmt.Errorf("requestID cannot be empty")
    }

    items := []model.StorageItem{
        {
            Name: "SpeechSynthesis",
            Store: func() error {
                return UploadBase64MP3ToS3(constant.APPCONFIG.S3AudioBucket, requestID, audioData)
            },
        },
        {
            Name: "LLM Response",
            Store: func() error {
                return UploadText(constant.APPCONFIG.S3LLMBucket, requestID, LLMResp)
            },
        },
        {
            Name: "SpeechMarks",
            Store: func() error {
                return UploadSpeechMarks(constant.APPCONFIG.S3SpeechMarkBucket, requestID, speechMarksData)
            },
        },
        {
            Name: "Cleaned Text",
            Store: func() error {
                return UploadText(constant.APPCONFIG.S3CleanedBucket, requestID, cleanedMarkdown)
            },
        },
    }

    var errs []string
    var successful []string
    
    for _, item := range items {
        if err := item.Store(); err != nil {
            errs = append(errs, fmt.Sprintf("%s: %v", item.Name, err))
            log.Error("Failed to store ", item.Name, ": ", err)
        } else {
            successful = append(successful, item.Name)
        }
    }

    if len(errs) > 0 {
        return fmt.Errorf("failed to store some items: %s", strings.Join(errs, "; "))
    }

    return nil
}