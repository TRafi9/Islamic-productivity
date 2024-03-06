FROM golang:1.21

WORKDIR /app
COPY . .

# RUN go install github.com/GoogleCloudPlatform/cloud-sql-proxy/v2@latest

RUN go build tpm

# RUN cloud-sql-proxy starlit-booster-408007:europe-west2:the-productive-muslim --credentials-file=/app/tpm_key.json

EXPOSE 8080

# cloud-sql-proxy starlit-booster-408007:europe-west2:the-productive-muslim --credentials-file=/app/tpm_key.json & sleep 5 && 

CMD ["sh", "-c", "./tpm"]