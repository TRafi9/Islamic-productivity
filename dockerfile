
## KUBERNETES DOCKER FILE VERSION
FROM golang:1.21

WORKDIR /app
COPY . .

RUN go build tpm

EXPOSE 8080


CMD ["sh", "-c", "./tpm"]

## LOCAL TESTING DOCKER FILE

# FROM golang:1.21

# WORKDIR /app
# COPY . .

# RUN go install github.com/GoogleCloudPlatform/cloud-sql-proxy/v2@latest

# RUN go build tpm

# # RUN cloud-sql-proxy starlit-booster-408007:europe-west2:the-productive-muslim --credentials-file=/app/tpm_key.json

# EXPOSE 8080

# CMD ["sh", "-c", "cloud-sql-proxy starlit-booster-408007:europe-west2:the-productive-muslim --credentials-file=/app/tpm_key.json & sleep 5 && ./tpm"]