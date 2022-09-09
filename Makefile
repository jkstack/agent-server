.PHONY: all

OUTDIR=release

VERSION=1.0.0
TIMESTAMP=`date +%s`

BRANCH=`git rev-parse --abbrev-ref HEAD`
HASH=`git log -n1 --pretty=format:%h`
REVERSION=`git log --oneline|wc -l|tr -d ' '`
BUILD_TIME=`date +'%Y-%m-%d %H:%M:%S'`
LDFLAGS="-X 'main.gitBranch=$(BRANCH)' \
-X 'main.gitHash=$(HASH)' \
-X 'main.gitReversion=$(REVERSION)' \
-X 'main.buildTime=$(BUILD_TIME)' \
-X 'main.version=$(VERSION)'"

all:
	go mod vendor
	rm -fr $(OUTDIR)/$(VERSION)
	mkdir -p $(OUTDIR)/$(VERSION)/opt/agent-server/bin \
		$(OUTDIR)/$(VERSION)/opt/agent-server/conf
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod vendor -ldflags $(LDFLAGS) \
		-o $(OUTDIR)/$(VERSION)/opt/agent-server/bin/agent-server main.go
	cp conf/server.conf $(OUTDIR)/$(VERSION)/opt/agent-server/conf/server.conf
	echo $(VERSION) > $(OUTDIR)/$(VERSION)/opt/agent-server/.version
	cd $(OUTDIR)/$(VERSION) && fakeroot tar -czvf agent-server_$(VERSION).tar.gz \
		--warning=no-file-changed opt
	rm -fr $(OUTDIR)/$(VERSION)/opt
	cp CHANGELOG.md $(OUTDIR)/CHANGELOG.md
clean:
	rm -fr $(OUTDIR)
version:
	@echo $(VERSION)
distclean: clean
	rm -fr vendor