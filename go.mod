module github.com/NorskHelsenett/terraform-provider-tss

require (
	github.com/danhale-git/tss-sdk-go v1.0.1
	github.com/hashicorp/terraform v0.12.14
	github.com/hashicorp/terraform-plugin-docs v0.5.0 // indirect
)

// replace github.com/thycotic/tss-sdk-go => ../tss-sdk-go

go 1.13
