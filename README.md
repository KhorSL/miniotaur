# Miniotaur

Simple object store service using [Minio](https://min.io/)

## Setup and start

```
cp .env_example .env
make start
```

### Accessing Minio

Credentials specified in `.env` file. See `.env_example` for example.

```
http://localhost:9001
```

### Create Buckets with AWS CLI Commands

```
# Access Key = MINIO_ROOT_USER
# Access Secret = MINIO_ROOT_PASSWORD
aws configure --profile minio
aws s3 mb s3://peanuts-bucket --endpoint-url http://localhost:9000 --profile=minio
aws s3 mb s3://strawberries-bucket --endpoint-url http://localhost:9000 --profile=minio
```

## Routes

```
GET /health
GET /api/v1/bucket
GET /api/v1/object/{bucket}/{key}
PUT /api/v1/object/{bucket}/{key}
```

## Examples

### Upload

```
curl -i -X PUT -H "Content-Type: multipart/form-data" -F "file=@./docs/test.json" http://localhost:8080/api/v1/object/strawberries-bucket/strawberry.json
```

### Upload Response

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Vary: Origin
Date: Thu, 29 Jun 2023 00:42:22 GMT
Content-Length: 103

{"bucket":"strawberries-bucket","etag":"\"3764d2e23e713c0b19726bcf05b781f6\"","key":"strawberry.json"}
```

### Get

```
curl -i -X GET http://localhost:8080/api/v1/object/strawberries-bucket/strawberry.json
```

### Get Response

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Vary: Origin
Date: Thu, 29 Jun 2023 00:45:45 GMT
Content-Length: 82

{"object":"ewogICJtb29kIjogInN0cmF3YmVycnkiLAogICJzdGF0dXMiOiAiZXhjZWxsZW50Igp9"}
```

#### Get Response decoded

```
$ echo ewogICJtb29kIjogInN0cmF3YmVycnkiLAogICJzdGF0dXMiOiAiZXhjZWxsZW50Igp9 | base64 --decode

{
  "mood": "strawberry",
  "status": "excellent"
}
```

## Alternative Minio image version to consider

- quay.io/minio/minio:RELEASE.2022-02-18T01-50-10Z

## References

- https://tutorialedge.net/golang/go-file-upload-tutorial/
- https://github.com/javiersoto15/skeleton-tutorial/tree/master
- https://aws.github.io/aws-sdk-go-v2/docs/handling-errors/
- https://itnext.io/structuring-a-production-grade-rest-api-in-golang-c0229b3feedc
