.PHONY: release
release:
	chmod +x ./scripts/release.sh && ./scripts/release.sh $(VERSION)

.PHONY: test
test:
	chmod +x ./scripts/test.sh && ./scripts/test.sh
