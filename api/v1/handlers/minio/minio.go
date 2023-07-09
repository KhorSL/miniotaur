package minio

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
	"github.com/khorsl/miniotaur/common/constants"
	"github.com/khorsl/miniotaur/common/log/logger"
	"github.com/khorsl/miniotaur/internal/services/minio"
	"github.com/rs/zerolog/log"
)

type MinioHandler struct {
	svc *minio.MinioSvc
}

func NewMinioHandler(svc *minio.MinioSvc) *MinioHandler {
	return &MinioHandler{
		svc: svc,
	}
}

func (mh *MinioHandler) ListBuckets(w http.ResponseWriter, r *http.Request) {
	logger := getLogger(os.Getenv("DEFAULT_LOGGER_TYPE"))

	buckets, err := mh.svc.ListBuckets(logger)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	result := make(map[string][]string)
	result["buckets"] = buckets

	render.JSON(w, r, result)
}

func (mh *MinioHandler) GetObject(w http.ResponseWriter, r *http.Request) {
	logger := getLogger(os.Getenv("DEFAULT_LOGGER_TYPE"))

	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "key")

	obj, err := mh.svc.GetObject(bucket, key, logger)
	if err != nil {
		//TODO
		var nsb *types.NoSuchBucket
		if errors.As(err, &nsb) {
			http.Error(w, "No such bucket", 404)
			return
		}

		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			http.Error(w, "No such key", 404)
			return
		}

		http.Error(w, http.StatusText(500), 500)
		return
	}

	result := make(map[string][]byte)
	result["object"] = obj

	render.JSON(w, r, result)
}

func (mh *MinioHandler) UploadObject(w http.ResponseWriter, r *http.Request) {
	logger := getLogger(os.Getenv("DEFAULT_LOGGER_TYPE"))

	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "key")

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Error("Error retrieving the file", map[string]interface{}{
			"error": err,
		})
		http.Error(w, err.Error(), 400)
		return
	}
	defer file.Close()

	logger.Info("Uploading file", map[string]interface{}{
		"filename":    handler.Filename,
		"file_size":   handler.Size,
		"mime_header": handler.Header,
	})

	etag, err := mh.svc.UploadObject(bucket, key, file, logger)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	logger.Info("File uploaded", map[string]interface{}{
		"filename":    handler.Filename,
		"file_size":   handler.Size,
		"mime_header": handler.Header,
	})

	result := make(map[string]string)
	result["etag"] = etag
	result["bucket"] = bucket
	result["key"] = key

	render.JSON(w, r, result)
}

func getLogger(loggerType string) logger.Logger {
	uid, err := uuid.NewV4()
	if err != nil {
		log.Error().Err(err).Msg("Unable to generate UUID")
	}

	ctx := context.TODO()
	ctx = context.WithValue(ctx, constants.LoggerCommonFields, map[string]interface{}{
		"correlation_id": uid.String(),
	})

	return logger.NewLoggerWrapper(loggerType, ctx)
}
