module github.com/NorskHelsenett/terraform-provider-tss

require (
	

	github.com/vidarno/tss-sdk-go/v2 v2.0.2
	github.com/hashicorp/terraform-plugin-docs v0.11.0 // indirect
	github.com/hashicorp/terraform-plugin-log v0.7.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.24.0
)

replace github.com/NorskHelsenett/terraform-provider-tss/tss => ./tss

go 1.13
