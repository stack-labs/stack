package template

var (
	StackConfig = `stack:
  server:
    name: {{.FQDN}}
    address: :0
  registry:
    name: mdns
`
)
