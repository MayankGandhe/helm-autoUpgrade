FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o helm-upgrade main.go

FROM alpine:3.19
RUN apk add --no-cache curl bash
# Install Helm
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 && \
    chmod 700 get_helm.sh && ./get_helm.sh
WORKDIR /app
COPY --from=builder /app/helm-upgrade .
ENTRYPOINT ["/app/helm-upgrade"]
