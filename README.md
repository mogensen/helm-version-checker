# helm-version-checker Helm Charts

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)![Release Charts](https://github.com/mogensen/helm-version-checker/workflows/Release%20Charts/badge.svg)

## Usage

To install using this Helm chart, first install Helm using Helm's [documentation](https://helm.sh/docs/).

Validate your Helm cli version

```bash
helm version
```

Use the repo as follows in Helm v3:

```bash
# Add repo to local helm setup
helm repo add helm-version-checker https://mogensen.github.io/helm-version-checker

# List charts and versions in the repo.
helm search repo helm-version-checker
```

## Contributing

The source code of the `helm-version-checker` [Helm](https://helm.sh) chart can be found on Github:
<https://github.com/mogensen/helm-version-checker/tree/main/deploy/charts/helm-version-checker>

## License

<!-- Keep full URL links to repo files because this README syncs from main to gh-pages.  -->
[Apache 2.0 License](https://github.com/mogensen/helm-version-checker/blob/main/LICENSE).

## Helm charts build status

![Release Charts](https://github.com/mogensen/helm-version-checker/workflows/Release%20Charts/badge.svg?branch=main)
