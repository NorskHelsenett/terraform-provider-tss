module github.com/thycotic/terraform-provider-tss

require (
	github.com/hashicorp/terraform v0.12.14
	github.com/danhale-git/tss-sdk-go v1.0.1
)

// replace github.com/thycotic/tss-sdk-go => ../tss-sdk-go

go 1.13
