#!/bin/bash


if ! helm version -c --short | grep -E "v3." >/dev/null; then
    echo "Helm v3 is needed!"
    exit 1
fi

helm template helm-version-checker deploy/charts/helm-version-checker --no-hooks --set image.pullPolicy=Always  \
    --set ingress.enabled=true  \
    --set 'ingress.hosts[0].host=helm-version-checker.localtest.me'  \
    --set 'ingress.hosts[0].paths[0].path="/"' \
    --set image.pullPolicy=IfNotPresent \
    | grep -vi "managed-by" \
    | grep -vi chart \
    | grep -v "# Source" > deploy/yaml/deploy.yaml

helm template helm-version-checker deploy/charts/helm-version-checker --no-hooks -s templates/grafana-dashboard-cm.yaml --set grafanaDashboard.enabled=true  \
    | grep -vi "managed-by" \
    | grep -vi chart \
    | grep -v "# Source" > deploy/yaml/grafana-dashboard-cm.yaml

helm template helm-version-checker deploy/charts/helm-version-checker --no-hooks -s templates/servicemonitor.yaml \
    --set serviceMonitor.enabled=true  \
    --set serviceMonitor.additionalLabels.release=prometheus  \
    | grep -vi "managed-by" \
    | grep -vi chart \
    | grep -v "# Source" > deploy/yaml/servicemonitor.yaml
