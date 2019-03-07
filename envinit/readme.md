## envinit

#### usage

```bash
GO = go
GOGET = ${GO} get
BUILD = ${CROSS} ${GO} build ${GOFLAGS}
BRANCH = ${shell basenae `git symbolic-ref HEAD`}
COMMIT = ${shell git rev-parse HEAD}
VERSION_NAMESPACE = github.com/go-infra/envinit
BUILD_VERSION = ${BUILD} -ldflags "-X $(VERSION_NAMESPACE).repository=$(REPO) -X $(VERSION_NAMESPACE).branch=$(BRANCH) -X $(VERSION_NAMESPACE).commit=$(COMMIT)"
