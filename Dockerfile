FROM python:3.12.10-alpine3.21

RUN apk add --no-cache curl bash openssl

RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 && \
    chmod 700 get_helm.sh && ./get_helm.sh

COPY . .

RUN pip3 install -r requirements.txt

ENTRYPOINT [ "uvicorn", "app:app", "--reload" , "--host", "0.0.0.0" ]