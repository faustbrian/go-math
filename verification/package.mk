.PHONY: api architecture docs provenance

api:
	./scripts/check-api-compat.sh

architecture:
	./scripts/check-architecture.sh

docs:
	./scripts/check-docs.sh

provenance:
	./scripts/check-provenance.sh
