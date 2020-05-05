package svd

var preambleTemplateText = `
package {{.Package}} 

{{/* Emit the type of each register */}}
{{range .Peripheral}}
{{range .Register}}
{{if ne .TypeName ""}}
type {{printf "%sDef" .Name}} volatile.Register32
{{end}} {{/*closes if*/}}
{{end}} {{/*closes registers*/}}
{{end}} {{/*closes peripherals*/}}
`

var deviceTemplateText = `


{{/* Emit the struct for each peripheral */}}
{{range .Peripheral}}
type {{printf "%sDef" .TypeName}} struct {
	{{- range .Register}}
		{{if ne .TypeName ""}} 
		{{.Name}} {{printf "%sDef" .Name}}   // 0x{{printf "%x" .AddressOffset.Get }}
		{{else}}
		{{.Name}} volatile.Register32   // 0x{{printf "%x" .AddressOffset.Get }}
		{{end}} {{/*closes if*/}}
	{{end}} {{/* closesregisters */}} 
} {{/* closes struct of peripheral */}}
var {{.TypeName}} {{printf "%sDef" .TypeName}}
{{end}} {{/* end of peripherals */}}
`

var bitFieldDeclTemplateText = `
{{if .CanRead}}
{{if eq .BitRange.Width 1}}
func (a *{{printf "%sDef" .RegName}}) {{printf "%sIsSet" .Name}}() bool {
	b:=BitField{ msb:{{ .BitRange.Msb}}, lsb:{{ .BitRange.Lsb}},reg:a}
	return b.HasBit()
}
{{else}}
func (a *{{printf "%sDef" .RegName}}) {{.Name}} () uint32 {
	b:=BitField{ msb:{{ .BitRange.Msb}}, lsb:{{ .BitRange.Lsb }} ,reg:a}
	return b.Get()
}
{{end}} {{/*end of the bit width is 1*/}}
{{end}} {{/*end of can read */}}

{{if .CanWrite}}
{{if eq .BitRange.Width 1}} 
func (a *{{printf "%sDef" .RegName}}) {{printf "%sSet" .Name}}() {
	b:=BitField{ msb:{{ .BitRange.Msb}}, lsb:{{ .BitRange.Lsb}},a}
	b.Set()
}
func (a *{{printf "%sDef" .RegName}}) {{printf "%sClear" .Name}}() {
	b:=BitField{ msb:{{ .BitRange.Msb}}, lsb:{{ .BitRange.Lsb}},a}
	b.Clear()
}
{{else}}
func (a *{{printf "%sDef" .RegName}}) {{printf "Set%s" .Name}}(v uint32) {
	b:=BitField{ msb:{{ .BitRange.Msb}}, lsb:{{ .BitRange.Lsb}},a}
	b.Set(v)
}
{{end}} {{/*closes if */}}

{{range .EnumeratedValue}}
{{if .Field.CanRead}}
func (a *{{printf "%sDef" .Field.RegName}}) {{.Name}}() bool {
	return a.Get()=={{ .Value.Get }}
}
{{end}} {{/* closes if */}}
{{if .Field.CanWrite}}
func (a *{{printf "%sDef" .Field.RegName }}) {{ printf "Set%s" .Name }}() bool {
	return a.Set({{ .Value.Get }})
}
{{end}} {{/* closes if */}}
{{end}} {{/*closes enumerated values*/}}
{{end}} {{/*closes Can write */}}
`
