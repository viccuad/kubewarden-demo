.PHONY: sigstore psp
sigstore:
	@clear
	@go run . --sigstore-demo

psp:
	@clear
	@go run . --psp-demo
