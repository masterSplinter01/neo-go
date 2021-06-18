package smartcontract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/nspcc-dev/neo-go/internal/random"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract/manifest"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/urfave/cli"
)

const srcTmpl = `package {{.PackageName}}

import (
	"github.com/nspcc-dev/neo-go/pkg/rpc/client"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract"
	"github.com/nspcc-dev/neo-go/pkg/util"
)

var contractHash = {{ printf "%#v" .Hash }}

// Client is a wrapper over RPC-client mirroring methods of smartcontract.
type Client client.Client
{{range $m := .Methods}}
// {{.Name}} {{.Comment}}
func (c *Client) {{.Name}}({{range $index, $arg := .Arguments -}}
	{{- if ne $index 0}}, {{end}}
		{{- .Name}} {{scTypeToGo .Type}}
	{{- end}}) {{if .Return }}({{ printType .Default | print}}, error){{else}}error{{end}} {
	args := make([]smartcontract.Parameter, {{ len .Arguments }})
	{{range $index, $arg := .Arguments -}}
	args[{{$index}}] = smartcontract.Parameter{Type: {{ scType $arg.Type }}, Value: {{ scName $arg.Type $arg.Name -}} }
	{{end}}
	result, err := (*client.Client)(c).InvokeFunction(contractHash, "
		{{- lowerFirst .Name }}", args, nil)
	if err != nil {
		return {{if .Return }}{{printValue .Default | print}}, {{end}}err
	}

	err = client.GetInvocationError(result)
	if err != nil {
		return {{if .Return }}{{.Default}}, {{end}}err
	}

	return client.TopIntFromStack(result.Stack)
}
{{end}}`

func printValue(v interface{}) string {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map, reflect.Interface, reflect.Slice:
		if rv.IsNil() {
			return "nil"
		}
	case reflect.String:
		return fmt.Sprintf("%q", v)
	}
	return fmt.Sprintf("%#v", v)
}

func printType(v interface{}) string {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map, reflect.Interface, reflect.Slice:
		if rv.IsNil() {
			return "interface{}"
		}
	}
	return fmt.Sprintf("%#T", v)
}

func scType(s smartcontract.ParamType) string {
	switch s {
	case smartcontract.IntegerType:
		return "smartcontract.IntegerType"
	case smartcontract.Hash160Type:
		return "smartcontract.Hash160Type"
	case smartcontract.VoidType:
		return ""
	default:
		return "smartcontract.AnyType"
	}
}

func scName(typ smartcontract.ParamType, name string) string {
	switch typ {
	case smartcontract.Hash160Type, smartcontract.Hash256Type:
		return name + ".BytesBE()"
	default:
		return name
	}
}

func scTypeToGo(typ smartcontract.ParamType) (string, interface{}) {
	switch typ {
	case smartcontract.AnyType:
		return "interface{}", interface{}(nil)
	case smartcontract.BoolType:
		return "bool", false
	case smartcontract.IntegerType:
		return "int64", 0
	case smartcontract.ByteArrayType, smartcontract.SignatureType, smartcontract.PublicKeyType:
		return "[]byte", []byte(nil)
	case smartcontract.StringType:
		return "string", `""`
	case smartcontract.Hash160Type:
		return "util.Uint160", util.Uint160{}
	case smartcontract.Hash256Type:
		return "util.Uint256", util.Uint256{}
	case smartcontract.ArrayType:
		return "[]interface{}", []interface{}(nil)
	case smartcontract.MapType:
		return "map[string]interface{}", map[string]interface{}(nil)
	case smartcontract.VoidType:
		return "", nil
	default:
		panic("unexpected")
	}
}

func upperFirst(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}
func lowerFirst(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

func Generate(arg interface{}) (string, error) {
	fm := template.FuncMap{
		"lowerFirst": lowerFirst,
		"scType":     scType,
		"scName":     scName,
		"scTypeToGo": func(s smartcontract.ParamType) string {
			typ, _ := scTypeToGo(s)
			return typ
		},
		"printValue": printValue,
		"printType":  printType,
	}
	tmp := template.New("test").Funcs(fm)
	tmp, err := tmp.Parse(srcTmpl)
	if err != nil {
		return "", err
	}
	b := bytes.NewBuffer(nil)
	if err := tmp.Execute(b, arg); err != nil {
		return "", err
	}
	return b.String(), nil
}

type (
	contractTmpl struct {
		PackageName string
		Hash        util.Uint160
		Methods     []methodTmpl
	}

	methodTmpl struct {
		Name      string
		Comment   string
		Arguments []manifest.Parameter
		Return    string
		Default   interface{}
	}
)

// contractGenerateWrapper deploys contract.
func contractGenerateWrapper(ctx *cli.Context) error {
	manifestFile := ctx.String("manifest")
	if len(manifestFile) == 0 {
		return cli.NewExitError(errNoManifestFile, 1)
	}

	manifestBytes, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("failed to read manifest file: %w", err), 1)
	}

	m := &manifest.Manifest{}
	err = json.Unmarshal(manifestBytes, m)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("failed to restore manifest file: %w", err), 1)
	}

	ctr := contractTmpl{
		PackageName: ctx.String("package"),
		Hash:        random.Uint160(),
	}

	for _, m := range m.ABI.Methods {
		if m.Name[0] == '_' {
			continue
		}
		typ, def := scTypeToGo(m.ReturnType)
		mtd := methodTmpl{
			Name:      upperFirst(m.Name),
			Return:    typ,
			Default:   def,
			Comment:   fmt.Sprintf("invokes `%s` method of contract.", m.Name),
			Arguments: m.Parameters,
		}
		ctr.Methods = append(ctr.Methods, mtd)
	}

	s, err := Generate(ctr)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("error during generation: %w", err), 1)
	}

	err = ioutil.WriteFile(ctx.String("out"), []byte(s), os.ModePerm)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("error during write: %w", err), 1)
	}
	return nil
}
