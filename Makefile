.PHONY: sigstore psp fill-empty
sigstore:
	@clear
	@go run . --sigstore-demo

psp:
	@clear
	@go run . --psp-demo

fill-empty:
	@cat empty.md && read x
