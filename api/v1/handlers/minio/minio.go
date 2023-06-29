package minio

import (
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/khorsl/minio_tutorial/internal/services/minio"
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
	buckets, err := mh.svc.ListBuckets()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	result := make(map[string][]string)
	result["buckets"] = buckets

	render.JSON(w, r, result)
}

func (mh *MinioHandler) GetObject(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "key")

	obj, err := mh.svc.GetObject(bucket, key)
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
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "key")

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving the file: %v", err)
		http.Error(w, err.Error(), 400)
		return
	}
	defer file.Close()

	log.Printf("Uploading file: %+v, File Size: %+v, MIME Header: %+v",
		handler.Filename,
		handler.Size,
		handler.Header)

	etag, err := mh.svc.UploadObject(bucket, key, file)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	result := make(map[string]string)
	result["etag"] = etag
	result["bucket"] = bucket
	result["key"] = key

	render.JSON(w, r, result)
}
