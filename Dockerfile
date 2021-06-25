ARG GLOBAL_KUBE_VERSION="v1.19.0"
ARG GLOBAL_HELM_VERSION="v3.5.4"


FROM golang:1.16-alpine  AS build-env
ARG GLOBAL_KUBE_VERSION
ARG GLOBAL_HELM_VERSION

ENV KUBE_VERSION=$GLOBAL_KUBE_VERSION
ENV HELM_VERSION=$GLOBAL_HELM_VERSION

RUN apk add --update --no-cache ca-certificates git openssh ruby curl tar gzip make bash


RUN curl --retry 5 -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl


RUN curl --retry 5 -Lk https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz | tar zxv -C /tmp
RUN mv /tmp/linux-amd64/helm /usr/local/bin/helm && rm -rf /tmp/linux-amd64
RUN chmod +x /usr/local/bin/helm


RUN helm plugin install https://github.com/fabmation-gmbh/helm-whatup

# Dependencies
WORKDIR /build
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY cmd cmd/
COPY pkg pkg/

RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o /app/helm-version-checker ./cmd/

# Build runtime container
FROM golang:1.16-alpine
LABEL description="Helm Chart version monitoring utility for watching updated and deprecated helm releases and reporting the result as metrics."
WORKDIR /app
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Setup user
RUN addgroup -g 35992 helm-version-checker && adduser -u 35992 -G helm-version-checker -D helm-version-checker -h  /home/helm-version-checker
USER helm-version-checker
ENV HOME=/home/helm-version-checker

COPY --from=build-env /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=build-env /usr/local/bin/helm /usr/local/bin/helm
COPY --from=build-env --chown=35992:35992 /root/.local /home/helm-version-checker/.local

# Needs a generic method for this
RUN helm repo add stable https://charts.helm.sh/stable
RUN helm repo add cert-checker https://mogensen.github.io/cert-checker
RUN helm repo update

COPY --from=build-env --chown=helm-version-checker /app/helm-version-checker /app/helm-version-checker

CMD ["/app/helm-version-checker"]
