package handlers

import (
	"github.com/go-chi/chi/v5"
	mh "github.com/khorsl/minio_tutorial/api/v1/handlers/minio"
	ms "github.com/khorsl/minio_tutorial/internal/services/minio"
)

func Routes() *chi.Mux {
	minioSvc := ms.NewMinioSvc()
	minioHandler := mh.NewMinioHandler(minioSvc)

	router := chi.NewRouter()

	// Bucket
	router.Get("/bucket", minioHandler.ListBuckets)

	// Object
	router.Get("/object/{bucket}/{key}", minioHandler.GetObject)
	router.Put("/object/{bucket}/{key}", minioHandler.UploadObject)

	return router
}
