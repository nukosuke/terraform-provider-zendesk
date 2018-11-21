# build provider & `terraform init`
init: ../*.go ../zendesk/*.go
	cd ../ && go build -o examples/terraform-provider-zendesk
	terraform init

# `terraform plan`
plan: init
	terraform plan
