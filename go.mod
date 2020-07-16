module github.com/LuciferInLove/dynamic-sshmenu-aws

go 1.14

require (
	github.com/aws/aws-sdk-go v1.33.4
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/manifoldco/promptui v0.7.0
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/urfave/cli/v2 v2.2.0
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
)

replace github.com/manifoldco/promptui v0.7.0 => github.com/LuciferInLove/promptui v0.7.1-0.20200604215815-d8893f35f691
