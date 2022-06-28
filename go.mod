module github.com/NorskHelsenett/terraform-provider-tss

require (
	github.com/hashicorp/terraform-plugin-log v0.4.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.17.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/thycotic/tss-sdk-go v1.2.1
)

replace github.com/NorskHelsenett/terraform-provider-tss/tss => ./tss

go 1.13
