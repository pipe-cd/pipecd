package main

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"strconv"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

type FileParams struct {
	InputPath string
	Methods   []*Method
}

type Method struct {
	Name     string // The name of the RPC
	Resource string // APPLICAION,DEPLOYMENT,EVENT,PIPED,DEPLOYMENT_CHAIN,PROJECT,API_KEY,INSIGHT
	Action   string // GET,LIST,CREATE,UPDATE,DELETE
	Ignored  bool   // Whether ignore authorization or not
}

const (
	filePrefix               = "pkg/app/server/service/webservice"
	generatedFileNameSuffix  = ".pb.auth.go"
	protoFileExtention       = ".proto"
	methodOptionsRBACResouce = "rbac_resource"
	methodOptionsRBACPolicy  = "rbac_policy"
	methodOptionsAuth        = "auth"
)

func main() {
	protogen.Options{}.Run(func(p *protogen.Plugin) error {
		extTypes := new(protoregistry.Types)
		for _, f := range p.Files {
			if err := registerAllExtensions(extTypes, f.Desc); err != nil {
				return fmt.Errorf("registerAllExtensions error: %v", err)
			}

			if !f.Generate || !strings.Contains(f.GeneratedFilenamePrefix, filePrefix) {
				continue
			}

			methods := make([]*Method, 0, len(f.Services)*len(f.Services[0].Methods))
			for _, svc := range f.Services {
				ms, err := generateMethods(extTypes, svc.Methods)
				if err != nil {
					return fmt.Errorf("generateMethods error: %v", err)
				}
				methods = append(methods, ms...)
			}

			filename := fmt.Sprintf("%s%s", f.GeneratedFilenamePrefix, generatedFileNameSuffix)
			gf := p.NewGeneratedFile(filename, f.GoImportPath)

			sort.SliceStable(methods, func(i, j int) bool {
				return methods[i].Resource < methods[j].Resource
			})

			inputPath := fmt.Sprintf("%s%s", f.GeneratedFilenamePrefix, protoFileExtention)
			fp := &FileParams{
				InputPath: inputPath,
				Methods:   methods,
			}

			buf := bytes.Buffer{}
			t := template.Must(template.New("auth").Parse(fileTpl))
			if err := t.Execute(&buf, fp); err != nil {
				return fmt.Errorf("template execute error: %v", err)
			}
			gf.P(string(buf.Bytes()))
		}
		return nil
	})
}

// generateMethods generates the []*Method from []*protogen.Method for pasing template.
// The MessageOptions as provided by protoc does not know about dynamically created extensions,
// so they are left as unknown fields. We round-trip marshal and unmarshal the options
// with a dynamically created resolver that does know about extensions at runtime.
// https://github.com/golang/protobuf/issues/1260#issuecomment-751517894
func generateMethods(extTypes *protoregistry.Types, ms []*protogen.Method) ([]*Method, error) {
	ret := make([]*Method, 0, len(ms))
	for _, m := range ms {
		opts := m.Desc.Options().(*descriptorpb.MethodOptions)
		raw, err := proto.Marshal(opts)
		if err != nil {
			return nil, err
		}

		opts.Reset()
		err = proto.UnmarshalOptions{Resolver: extTypes}.Unmarshal(raw, opts)
		if err != nil {
			return nil, err
		}

		method := &Method{Name: m.GoName}
		opts.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
			if !fd.IsExtension() {
				return true
			}

			var value string
			if v.String() != "" {
				value = strings.SplitN(v.String(), ":", 2)[1]
			}

			switch fd.Name() {
			case methodOptionsRBACResouce:
				method.Resource = value
			case methodOptionsRBACPolicy:
				method.Action = value
			case methodOptionsAuth:
				if v, err := strconv.ParseBool(value); err == nil {
					method.Ignored = v
				}
			}
			return true
		})

		if method.Ignored || (method.Resource != "" && method.Action != "") {
			ret = append(ret, method)
		}
	}
	return ret, nil
}

// registerAllExtensions recursively registers all extensions into the provided protoregistry.Types,
// starting with the protoreflect.FileDescriptor and recursing into its MessageDescriptors,
// their nested MessageDescriptors, and so on.
//
// This leverages the fact that both protoreflect.FileDescriptor and protoreflect.MessageDescriptor
// have identical Messages() and Extensions() functions in order to recurse through a single function.
// https://github.com/golang/protobuf/issues/1260#issuecomment-751517894
func registerAllExtensions(extTypes *protoregistry.Types, descs interface {
	Messages() protoreflect.MessageDescriptors
	Extensions() protoreflect.ExtensionDescriptors
}) error {
	mds := descs.Messages()
	for i := 0; i < mds.Len(); i++ {
		registerAllExtensions(extTypes, mds.Get(i))
	}
	xds := descs.Extensions()
	for i := 0; i < xds.Len(); i++ {
		if err := extTypes.RegisterExtension(dynamicpb.NewExtensionType(xds.Get(i))); err != nil {
			return err
		}
	}
	return nil
}
