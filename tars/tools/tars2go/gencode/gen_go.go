package gencode

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/TarsCloud/TarsGo/tars/tools/tars2go/ast"
	"github.com/TarsCloud/TarsGo/tars/tools/tars2go/options"
	"github.com/TarsCloud/TarsGo/tars/tools/tars2go/parse"
	"github.com/TarsCloud/TarsGo/tars/tools/tars2go/token"
	"github.com/TarsCloud/TarsGo/tars/tools/tars2go/utils"
	"github.com/TarsCloud/TarsGo/tars/tools/tars2go/version"
)

var (
	fileMap sync.Map
)

// GenGo record go code information.
type GenGo struct {
	opt      *options.Options
	code     bytes.Buffer
	vc       int // var count. Used to generate unique variable names
	filepath string
	prefix   string
	module   *ast.ModuleInfo

	// proto file name(not include .tars)
	ProtoName string
}

// NewGenGo build up a new filepath
func NewGenGo(opt *options.Options, filepath string) *GenGo {
	if opt.Outdir != "" && !strings.HasSuffix(opt.Outdir, "/") {
		opt.Outdir += "/"
	}

	return &GenGo{opt: opt, filepath: filepath, prefix: opt.Outdir, ProtoName: utils.Path2ProtoName(filepath)}
}

func getShortTypeName(src string) string {
	vec := strings.Split(src, "::")
	return vec[len(vec)-1]
}

func errString(hasRet bool) string {
	var retStr string
	if hasRet {
		retStr = "return ret, err"
	} else {
		retStr = "return err"
	}
	return `if err != nil {
  ` + retStr + `
  }`
}

func genForHead(vc string) string {
	i := `i` + vc
	e := `e` + vc
	return ` for ` + i + `,` + e + ` := int32(0), length;` + i + `<` + e + `;` + i + `++ `
}

// Gen to parse file.
func (g *GenGo) Gen() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			// set exit code
			os.Exit(1)
		}
	}()

	g.module = parse.NewParse(g.opt, g.filepath, make([]string, 0))
	g.genAll()
}

func (g *GenGo) P(v ...string) {
	for _, x := range v {
		fmt.Fprint(&g.code, x)
	}
	fmt.Fprintln(&g.code)
}

func (g *GenGo) W(v ...string) {
	for _, x := range v {
		fmt.Fprint(&g.code, x)
	}
}

func (g *GenGo) genAll() {
	if _, ok := fileMap.Load(g.filepath); ok {
		// already compiled
		return
	}
	fileMap.Store(g.filepath, struct{}{})

	g.module.Rename(g.opt.ModuleUpper)
	g.genInclude(g.module.IncModule)

	g.code.Reset()
	g.genHead()
	g.genPackage()

	for _, v := range g.module.Enum {
		g.genEnum(&v)
	}

	g.genConst(g.module.Const)

	for _, v := range g.module.Struct {
		g.genStruct(&v)
	}
	if len(g.module.Enum) > 0 || len(g.module.Const) > 0 || len(g.module.Struct) > 0 {
		g.saveToSourceFile(utils.Path2ProtoName(g.filepath) + ".go")
	}

	for _, v := range g.module.Interface {
		g.genInterface(&v)
	}
}

func (g *GenGo) genErr(err string) {
	panic(err)
}

func (g *GenGo) saveToSourceFile(filename string) {
	var beauty []byte
	var err error
	prefix := g.prefix

	if !g.opt.E {
		beauty, err = format.Source(g.code.Bytes())
		if err != nil {
			if g.opt.Debug {
				fmt.Println("------------------")
				fmt.Println(string(g.code.Bytes()))
				fmt.Println("------------------")
			}
			g.genErr("go fmt fail. " + filename + " " + err.Error())
		}
	} else {
		beauty = g.code.Bytes()
	}

	if filename == "stdout" {
		fmt.Println(string(beauty))
	} else {
		var mkPath string
		if g.opt.ModuleCycle {
			mkPath = prefix + g.ProtoName + "/" + g.module.Name
		} else {
			mkPath = prefix + g.module.Name
		}
		err = os.MkdirAll(mkPath, 0766)
		if err != nil {
			g.genErr(err.Error())
		}

		err = os.WriteFile(mkPath+"/"+filename, beauty, 0666)
		if err != nil {
			g.genErr(err.Error())
		}
	}
}

func (g *GenGo) genVariableName(prefix, name string) string {
	if strings.HasPrefix(name, "(*") && strings.HasSuffix(name, ")") {
		return strings.Trim(name, "()")
	}
	return prefix + name
}

func (g *GenGo) genHead() {
	g.P("// Code generated by tars2go ", version.VERSION, ", DO NOT EDIT.")
	g.P("// This file was generated from ", filepath.Base(g.filepath))
	g.P("// Package ", g.module.Name, " comment")
}

func (g *GenGo) genPackage() {
	g.P("package ", g.module.Name)
	g.P()
	g.P("import (")
	g.P(strconv.Quote("fmt"))
	g.P()
	g.P(strconv.Quote(g.opt.TarsPath + "/protocol/codec"))

	mImports := make(map[string]bool)
	for _, st := range g.module.Struct {
		if g.opt.ModuleCycle {
			for k, v := range st.DependModuleWithJce {
				g.genStructImport(k, v, mImports)
			}
		} else {
			for k := range st.DependModule {
				g.genStructImport(k, "", mImports)
			}
		}
	}
	for path := range mImports {
		g.P(path)
	}
	g.P(")")
	g.P()
	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var _ = fmt.Errorf")
	g.P("var _ = codec.FromInt8")
}

func (g *GenGo) genStructImport(module string, protoName string, mImports map[string]bool) {
	var moduleStr string
	var jcePath string
	var moduleAlia string
	if g.opt.ModuleCycle {
		moduleStr = module[len(protoName)+1:]
		jcePath = protoName + "/"
		moduleAlia = module + " "
	} else {
		moduleStr = module
	}

	for _, p := range g.opt.Imports {
		if strings.HasSuffix(p, "/"+moduleStr) {
			mImports[strconv.Quote(p)] = true
			return
		}
	}

	if g.opt.ModuleUpper {
		moduleAlia = utils.UpperFirstLetter(moduleAlia)
	}

	// example:
	// TarsTest.tars, MyApp
	// gomod:
	// github.com/xxx/yyy/tars-protocol/MyApp
	// github.com/xxx/yyy/tars-protocol/TarsTest/MyApp
	//
	// gopath:
	// MyApp
	// TarsTest/MyApp
	var modulePath string
	if g.opt.Module != "" {
		mf := filepath.Clean(filepath.Join(g.opt.Module, g.prefix))
		if runtime.GOOS == "windows" {
			mf = strings.ReplaceAll(mf, string(os.PathSeparator), string('/'))
		}
		modulePath = fmt.Sprintf("%s/%s%s", mf, jcePath, moduleStr)
	} else {
		modulePath = fmt.Sprintf("%s%s", jcePath, moduleStr)
	}
	mImports[moduleAlia+strconv.Quote(modulePath)] = true
}

func (g *GenGo) genIFPackage(itf *ast.InterfaceInfo) {
	g.P("package " + g.module.Name)
	g.P()

	g.P("import (")
	g.P(strconv.Quote("bytes"))
	g.P(strconv.Quote("context"))
	g.P(strconv.Quote("encoding/json"))
	g.P(strconv.Quote("fmt"))
	//g.P(strconv.Quote("unsafe"))
	g.P()

	tarsPath := g.opt.TarsPath
	if g.opt.AddServant || !g.opt.WithoutTrace {
		g.P(strconv.Quote(tarsPath))
	}
	g.P(strconv.Quote(tarsPath + "/model"))
	g.P(strconv.Quote(tarsPath + "/protocol/res/requestf"))
	g.P(strconv.Quote(tarsPath + "/protocol/codec"))
	g.P(strconv.Quote(tarsPath + "/protocol/tup"))
	g.P(strconv.Quote(tarsPath + "/protocol/res/basef"))
	g.P(strconv.Quote(tarsPath + "/util/tools"))
	g.P(strconv.Quote(tarsPath + "/util/endpoint"))
	g.P(strconv.Quote(tarsPath + "/util/current"))
	if !g.opt.WithoutTrace {
		g.P("tarstrace ", strconv.Quote(tarsPath+"/util/trace"))
	}

	if g.opt.ModuleCycle {
		for k, v := range itf.DependModuleWithJce {
			g.genIFImport(k, v)
		}
	} else {
		for k := range itf.DependModule {
			g.genIFImport(k, "")
		}
	}

	g.P(")")
	g.P()
	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var (")
	g.P("	_ = fmt.Errorf")
	g.P("	_ = codec.FromInt8")
	//g.P("	_ = unsafe.Pointer(nil)")
	g.P("	_ = bytes.ErrTooLarge")
	g.P(")")
}

func (g *GenGo) genIFImport(module string, protoName string) {
	var moduleStr string
	var jcePath string
	var moduleAlia string
	if g.opt.ModuleCycle {
		moduleStr = module[len(protoName)+1:]
		jcePath = protoName + "/"
		moduleAlia = module + " "
	} else {
		moduleStr = module
	}
	for _, p := range g.opt.Imports {
		if strings.HasSuffix(p, "/"+moduleStr) {
			g.P(strconv.Quote(p))
			return
		}
	}

	if g.opt.ModuleUpper {
		moduleAlia = utils.UpperFirstLetter(moduleAlia)
	}

	// example:
	// TarsTest.tars, MyApp
	// gomod:
	// github.com/xxx/yyy/tars-protocol/MyApp
	// github.com/xxx/yyy/tars-protocol/TarsTest/MyApp
	//
	// gopath:
	// MyApp
	// TarsTest/MyApp
	var modulePath string
	if g.opt.Module != "" {
		mf := filepath.Clean(filepath.Join(g.opt.Module, g.prefix))
		if runtime.GOOS == "windows" {
			mf = strings.ReplaceAll(mf, string(os.PathSeparator), string('/'))
		}
		modulePath = fmt.Sprintf("%s/%s%s", mf, jcePath, moduleStr)
	} else {
		modulePath = fmt.Sprintf("%s%s", jcePath, moduleStr)
	}
	g.P(moduleAlia, strconv.Quote(modulePath))
}

func (g *GenGo) typeDef(mb *ast.StructMember) string {
	if mb.Default != "" {
		return mb.Default
	}
	switch mb.Type.Type {
	case token.TBool:
		return "false"
	case token.TShort, token.TInt, token.TLong, token.TFloat, token.TDouble:
		return "0"
	case token.TString:
		return strconv.Quote("")
	default:
		g.genErr("typeDef unknown type")
	}
	return ""
}

func (g *GenGo) genType(ty *ast.VarType) string {
	ret := ""
	switch ty.Type {
	case token.TBool:
		ret = "bool"
	case token.TInt:
		if ty.Unsigned {
			ret = "uint32"
		} else {
			ret = "int32"
		}
	case token.TShort:
		if ty.Unsigned {
			ret = "uint16"
		} else {
			ret = "int16"
		}
	case token.TByte:
		if ty.Unsigned {
			ret = "uint8"
		} else {
			ret = "int8"
		}
	case token.TLong:
		if ty.Unsigned {
			ret = "uint64"
		} else {
			ret = "int64"
		}
	case token.TFloat:
		ret = "float32"
	case token.TDouble:
		ret = "float64"
	case token.TString:
		ret = "string"
	case token.TVector:
		ret = "[]" + g.genType(ty.TypeK)
	case token.TMap:
		ret = "map[" + g.genType(ty.TypeK) + "]" + g.genType(ty.TypeV)
	case token.Name:
		ret = strings.Replace(ty.TypeSt, "::", ".", -1)
		vec := strings.Split(ty.TypeSt, "::")
		for i := range vec {
			if g.opt.ModuleUpper {
				vec[i] = utils.UpperFirstLetter(vec[i])
			} else {
				if i == (len(vec) - 1) {
					vec[i] = utils.UpperFirstLetter(vec[i])
				}
			}
		}
		ret = strings.Join(vec, ".")
	case token.TArray:
		ret = "[" + fmt.Sprintf("%v", ty.TypeL) + "]" + g.genType(ty.TypeK)
	default:
		g.genErr("Unknown Type " + token.Value(ty.Type))
	}
	return ret
}

func (g *GenGo) genStructDefine(st *ast.StructInfo) {
	g.P("// ", st.Name, " struct implement")
	g.P("type ", st.Name, " struct {")
	for _, v := range st.Mb {
		tarsTag := `tars:"` + v.OriginKey + `,tag:` + strconv.Itoa(int(v.Tag)) + `,require:` + strconv.FormatBool(v.Require) + `"`
		if g.opt.JsonOmitEmpty {
			g.P(v.Key, " ", g.genType(v.Type), " `json:\"", v.OriginKey, ",omitempty\" ", tarsTag, "`")
		} else {
			g.P(v.Key, " ", g.genType(v.Type), " `json:\"", v.OriginKey, "\" ", tarsTag, "`")
		}
	}
	g.P("}")
}

func (g *GenGo) genFunResetDefault(st *ast.StructInfo) {
	g.P("func (st *", st.Name, ") ResetDefault() {")
	for _, v := range st.Mb {
		if v.Type.CType == token.Struct {
			g.P("st.", v.Key, ".ResetDefault()")
		}
		if v.Default == "" {
			continue
		}
		g.P("st.", v.Key, " = ", v.Default)
	}
	g.P("}")
}

func (g *GenGo) genWriteSimpleList(mb *ast.StructMember, prefix string, hasRet bool) {
	tag := strconv.Itoa(int(mb.Tag))
	unsigned := "Int8"
	if mb.Type.TypeK.Unsigned {
		unsigned = "Uint8"
	}
	errStr := errString(hasRet)
	g.P("err = buf.WriteHead(codec.SimpleList, ", tag, ")")
	g.P(errStr)
	g.P("err = buf.WriteHead(codec.BYTE, 0)")
	g.P(errStr)
	g.P("err = buf.WriteInt32(int32(len(", g.genVariableName(prefix, mb.Key), ")), 0)")
	g.P(errStr)
	g.P("err = buf.WriteSlice", unsigned, "(", g.genVariableName(prefix, mb.Key), ")")
	g.P(errStr)
}

func (g *GenGo) genWriteVector(mb *ast.StructMember, prefix string, hasRet bool) {
	if !mb.Require {
		g.P("if len(", g.genVariableName(prefix, mb.Key), ") > 0 {")
		defer g.P("}")
	}
	// SimpleList
	if mb.Type.TypeK.Type == token.TByte && !mb.Type.TypeK.Unsigned {
		g.genWriteSimpleList(mb, prefix, hasRet)
		return
	}
	// LIST
	errStr := errString(hasRet)
	tag := strconv.Itoa(int(mb.Tag))
	g.P("err = buf.WriteHead(codec.LIST, ", tag, ")")
	g.P(errStr)
	g.P("err = buf.WriteInt32(int32(len(", g.genVariableName(prefix, mb.Key), ")), 0)")
	g.P(errStr)
	g.P("for _, v := range ", g.genVariableName(prefix, mb.Key), " {")
	// for _, v := range can nesting for _, v := range，does not conflict, support multidimensional arrays

	dummy := &ast.StructMember{
		Require: true,
		Type:    mb.Type.TypeK,
		Key:     "v",
	}
	g.genWriteVar(dummy, "", hasRet)

	g.P("}")
}

func (g *GenGo) genWriteArray(mb *ast.StructMember, prefix string, hasRet bool) {
	if !mb.Require {
		g.P("if len(", g.genVariableName(prefix, mb.Key), ") > 0 {")
		defer g.P("}")
	}

	// SimpleList
	if mb.Type.TypeK.Type == token.TByte && !mb.Type.TypeK.Unsigned {
		g.genWriteSimpleList(mb, prefix, hasRet)
		return
	}

	// LIST
	errStr := errString(hasRet)
	tag := strconv.Itoa(int(mb.Tag))
	g.P("err = buf.WriteHead(codec.LIST, ", tag, ")")
	g.P(errStr)
	g.P("err = buf.WriteInt32(int32(len(", g.genVariableName(prefix, mb.Key), ")), 0)")
	g.P(errStr)
	// for _, v := range can nesting for _, v := range，does not conflict, support multidimensional arrays
	g.P("for _, v := range ", g.genVariableName(prefix, mb.Key), " {")
	dummy := &ast.StructMember{
		Require: true,
		Type:    mb.Type.TypeK,
		Key:     "v",
	}
	g.genWriteVar(dummy, "", hasRet)
	g.P("}")
}

func (g *GenGo) genWriteMap(mb *ast.StructMember, prefix string, hasRet bool) {
	if !mb.Require {
		g.P("if len(", g.genVariableName(prefix, mb.Key), ") > 0 {")
		defer g.P("}")
	}

	tag := strconv.Itoa(int(mb.Tag))
	vc := strconv.Itoa(g.vc)
	g.vc++
	errStr := errString(hasRet)
	g.P("err = buf.WriteHead(codec.MAP, ", tag, ")")
	g.P(errStr)
	g.P("err = buf.WriteInt32(int32(len(", g.genVariableName(prefix, mb.Key), ")), 0)")
	g.P(errStr)
	g.P("for k", vc, ", v", vc, " := range ", g.genVariableName(prefix, mb.Key), " {")
	// for _, v := range can nesting for _, v := range，does not conflict, support multidimensional arrays

	dummy := &ast.StructMember{
		Require: true,
		Type:    mb.Type.TypeK,
		Key:     "k" + vc,
	}
	g.genWriteVar(dummy, "", hasRet)

	dummy = &ast.StructMember{
		Tag:     1,
		Require: true,
		Type:    mb.Type.TypeV,
		Key:     "v" + vc,
	}
	g.genWriteVar(dummy, "", hasRet)

	g.P("}")
}

func (g *GenGo) genWriteStruct(mb *ast.StructMember, prefix string, hasRet bool) {
	tag := strconv.Itoa(int(mb.Tag))
	g.P("err = ", prefix, mb.Key, ".WriteBlock(buf, ", tag, ")")
	g.P(errString(hasRet))
}

func (g *GenGo) genWriteVar(v *ast.StructMember, prefix string, hasRet bool) {
	switch v.Type.Type {
	case token.TVector:
		g.genWriteVector(v, prefix, hasRet)
	case token.TArray:
		g.genWriteArray(v, prefix, hasRet)
	case token.TMap:
		g.genWriteMap(v, prefix, hasRet)
	case token.Name:
		if v.Type.CType == token.Enum {
			// tkEnum enumeration processing
			tag := strconv.Itoa(int(v.Tag))
			g.P("err = buf.WriteInt32(int32(", g.genVariableName(prefix, v.Key), "), ", tag, ")")
			g.P(errString(hasRet))
		} else {
			g.genWriteStruct(v, prefix, hasRet)
		}
	default:
		if !v.Require {
			g.P("if ", g.genVariableName(prefix, v.Key), " != ", g.typeDef(v), " {")
			defer g.P("}")
		}
		tag := strconv.Itoa(int(v.Tag))
		g.P("err = buf.Write", utils.UpperFirstLetter(g.genType(v.Type)), "(", g.genVariableName(prefix, v.Key), ", ", tag, ")")
		g.P(errString(hasRet))
	}
}

func (g *GenGo) genFunWriteBlock(st *ast.StructInfo) {
	// WriteBlock function head
	g.P("// WriteBlock encode struct")
	g.P("func (st *", st.Name, ") WriteBlock(buf *codec.Buffer, tag byte) error {")
	g.P("var err error")
	g.P("err = buf.WriteHead(codec.StructBegin, tag)")
	g.P("if err != nil {")
	g.P("	return err")
	g.P("}")
	g.P()
	g.P("err = st.WriteTo(buf)")
	g.P("if err != nil {")
	g.P("	return err")
	g.P("}")
	g.P()
	g.P("err = buf.WriteHead(codec.StructEnd, 0)")
	g.P("if err != nil {")
	g.P("	return err")
	g.P("}")
	g.P("return nil")
	g.P("}")
}

func (g *GenGo) genFunWriteTo(st *ast.StructInfo) {
	g.P("// WriteTo encode struct to buffer")
	g.P("func (st *", st.Name, ") WriteTo(buf *codec.Buffer) (err error) {")
	for _, v := range st.Mb {
		g.genWriteVar(&v, "st.", false)
		g.P()
	}
	g.P("return err")
	g.P("}")
}

func (g *GenGo) genReadSimpleList(mb *ast.StructMember, prefix string, hasRet bool) {
	unsigned := "Int8"
	if mb.Type.TypeK.Unsigned {
		unsigned = "Uint8"
	}
	errStr := errString(hasRet)
	g.P("_, err = readBuf.SkipTo(codec.BYTE, 0, true)")
	g.P(errStr)
	g.P("err = readBuf.ReadInt32(&length, 0, true)")
	g.P(errStr)
	g.P("err = readBuf.ReadSlice", unsigned, `(&`, prefix, mb.Key, ", length, true)")
	g.P(errStr)
}

func (g *GenGo) genReadVector(mb *ast.StructMember, prefix string, hasRet bool) {
	// LIST
	errStr := errString(hasRet)
	tag := strconv.Itoa(int(mb.Tag))
	vc := strconv.Itoa(g.vc)
	g.vc++
	if mb.Require {
		g.P("_, ty, err = readBuf.SkipToNoCheck(", tag, ", true)")
		g.P(errStr)
	} else {
		g.P("have, ty, err = readBuf.SkipToNoCheck(", tag, ", false)")
		g.P(errStr)
		g.P("if have {")
		// 结束标记
		defer g.P("}")
	}

	g.P("if ty == codec.LIST {")
	g.P("err = readBuf.ReadInt32(&length, 0, true)")
	g.P(errStr)
	g.P(g.genVariableName(prefix, mb.Key), " = make(", g.genType(mb.Type), ", length)")
	g.P(genForHead(vc), "{")

	dummy := &ast.StructMember{
		Require: true,
		Type:    mb.Type.TypeK,
		Key:     mb.Key + "[i" + vc + "]",
	}
	g.genReadVar(dummy, prefix, hasRet)

	g.P("}")
	g.P("} else if ty == codec.SimpleList {")
	if mb.Type.TypeK.Type == token.TByte {
		g.genReadSimpleList(mb, prefix, hasRet)
	} else {
		g.P(`err = fmt.Errorf("not support SimpleList type")`)
		g.P(errStr)
	}
	g.P("} else {")
	g.P(`err = fmt.Errorf("require vector, but not")`)
	g.P(errStr)
	g.P("}")
}

func (g *GenGo) genReadArray(mb *ast.StructMember, prefix string, hasRet bool) {
	// LIST
	errStr := errString(hasRet)
	tag := strconv.Itoa(int(mb.Tag))
	vc := strconv.Itoa(g.vc)
	g.vc++

	g.P()
	if mb.Require {
		g.P("_, ty, err = readBuf.SkipToNoCheck(", tag, ", true)")
		g.P(errStr)
	} else {
		g.P("have, ty, err = readBuf.SkipToNoCheck(", tag, ", false)")
		g.P(errStr)
		g.P("if have {")
		// 结束标记
		defer g.P("}")
	}

	g.P("if ty == codec.LIST {")
	g.P("err = readBuf.ReadInt32(&length, 0, true)")
	g.P(errStr)
	g.P(genForHead(vc), "{")
	dummy := &ast.StructMember{
		Require: true,
		Type:    mb.Type.TypeK,
		Key:     mb.Key + "[i" + vc + "]",
	}
	g.genReadVar(dummy, prefix, hasRet)
	g.P("}")

	g.P("} else if ty == codec.SimpleList {")
	if mb.Type.TypeK.Type == token.TByte {
		g.genReadSimpleList(mb, prefix, hasRet)
	} else {
		g.P(`err = fmt.Errorf("not support SimpleList type")`)
		g.P(errStr)
	}
	g.P("} else {")
	g.P(`err = fmt.Errorf("require array, but not")`)
	g.P(errStr)
	g.P("}")
}

func (g *GenGo) genReadMap(mb *ast.StructMember, prefix string, hasRet bool) {
	tag := strconv.Itoa(int(mb.Tag))
	errStr := errString(hasRet)
	vc := strconv.Itoa(g.vc)
	g.vc++

	if mb.Require {
		g.P("_, err = readBuf.SkipTo(codec.MAP, ", tag, ", true)")
		g.P(errStr)
	} else {
		g.P("have, err = readBuf.SkipTo(codec.MAP, ", tag, ", false)")
		g.P(errStr)
		g.P("if have {")
		// 结束标记
		defer g.P("}")
	}

	g.P("err = readBuf.ReadInt32(&length, 0, true)")
	g.P(errStr)
	g.P(g.genVariableName(prefix, mb.Key), " = make(", g.genType(mb.Type), ")")
	g.P(genForHead(vc), "{")
	g.P("var k", vc, " ", g.genType(mb.Type.TypeK))
	g.P("var v", vc, " ", g.genType(mb.Type.TypeV))

	key := "k" + vc
	dummy := &ast.StructMember{
		Require: true,
		Type:    mb.Type.TypeK,
		Key:     key,
	}
	g.genReadVar(dummy, "", hasRet)

	val := "v" + vc
	dummy = &ast.StructMember{
		Tag:     1,
		Require: true,
		Type:    mb.Type.TypeV,
		Key:     val,
	}
	g.genReadVar(dummy, "", hasRet)
	g.P(prefix, mb.Key, "[", key, "] = ", val)
	g.P("}")
}

func (g *GenGo) genReadStruct(mb *ast.StructMember, prefix string, hasRet bool) {
	tag := strconv.Itoa(int(mb.Tag))
	require := strconv.FormatBool(mb.Require)
	g.P("err = ", prefix, mb.Key, ".ReadBlock(readBuf, ", tag, ", ", require, ")")
	g.P(errString(hasRet))
}

func (g *GenGo) genReadVar(v *ast.StructMember, prefix string, hasRet bool) {
	switch v.Type.Type {
	case token.TVector:
		g.genReadVector(v, prefix, hasRet)
	case token.TArray:
		g.genReadArray(v, prefix, hasRet)
	case token.TMap:
		g.genReadMap(v, prefix, hasRet)
	case token.Name:
		if v.Type.CType == token.Enum {
			tag := strconv.Itoa(int(v.Tag))
			require := strconv.FormatBool(v.Require)
			g.P("err = readBuf.ReadInt32((*int32)(&", prefix, v.Key, "),", tag, ", ", require, ")")
			g.P(errString(hasRet))
		} else {
			g.genReadStruct(v, prefix, hasRet)
		}
	default:
		tag := strconv.Itoa(int(v.Tag))
		require := strconv.FormatBool(v.Require)
		g.P("err = readBuf.Read", utils.UpperFirstLetter(g.genType(v.Type)), "(&", prefix, v.Key, ", ", tag, ", ", require, ")")
		g.P(errString(hasRet))
	}
}

func (g *GenGo) genFunReadFrom(st *ast.StructInfo) {
	g.P(`// ReadFrom reads  from readBuf and put into struct.
func (st *`, st.Name, `) ReadFrom(readBuf *codec.Reader) error {
	var (
		err error
		length int32
		have bool
		ty byte
	)
	st.ResetDefault()
`)

	for _, v := range st.Mb {
		g.genReadVar(&v, "st.", false)
		g.P()
	}

	g.P(`_ = err
			 _ = length
			 _ = have
			 _ = ty
			 return nil
			}`)
}

func (g *GenGo) genFunReadBlock(st *ast.StructInfo) {
	g.P(`// ReadBlock reads struct from the given tag , require or optional.
func (st *`, st.Name, `) ReadBlock(readBuf *codec.Reader, tag byte, require bool) error {
	var (
		err error
		have bool
	)
	st.ResetDefault()

	have, err = readBuf.SkipTo(codec.StructBegin, tag, require)
	if err != nil {
		return err
	}
	if !have {
		if require {
			return fmt.Errorf("require `, st.Name, `, but not exist. tag %d", tag)
		}
		return nil
	}

  	err = st.ReadFrom(readBuf)
  	if err != nil {
		return err
	}

	err = readBuf.SkipToStructEnd()
	if err != nil {
		return err
	}
	_ = have
	return nil
}`)
}

func (g *GenGo) genStruct(st *ast.StructInfo) {
	g.vc = 0
	st.Rename()

	g.genStructDefine(st)
	g.genFunResetDefault(st)

	g.genFunReadFrom(st)
	g.genFunReadBlock(st)

	g.genFunWriteTo(st)
	g.genFunWriteBlock(st)
}

func (g *GenGo) makeEnumName(en *ast.EnumInfo, mb *ast.EnumMember) string {
	return utils.UpperFirstLetter(en.Name) + "_" + utils.UpperFirstLetter(mb.Key)
}

func (g *GenGo) genEnum(en *ast.EnumInfo) {
	if len(en.Mb) == 0 {
		return
	}

	en.Rename()

	g.P("//go:generate stringer -type " + en.Name + " -trimprefix " + en.Name + "_ -output " + strings.ToLower(en.Name) + "_string.go")
	g.P("type ", en.Name, " int32")
	g.P("const (")
	var it int32
	for _, v := range en.Mb {
		if v.Type == 0 {
			//use value
			g.P(g.makeEnumName(en, &v), " ", en.Name, " = ", strconv.Itoa(int(v.Value)))
			it = v.Value + 1
		} else if v.Type == 1 {
			// use name
			find := false
			for _, ref := range en.Mb {
				if ref.Key == v.Name {
					find = true
					g.P(g.makeEnumName(en, &v), " ", en.Name, " = ", g.makeEnumName(en, &ref))
					it = ref.Value + 1
					break
				}
				if ref.Key == v.Key {
					break
				}
			}
			if !find {
				g.genErr(v.Name + " not define before use.")
			}
		} else {
			// use auto add
			g.P(g.makeEnumName(en, &v), " ", en.Name, " = ", strconv.Itoa(int(it)))
			it++
		}
	}
	g.P(")")
}

func (g *GenGo) genConst(cst []ast.ConstInfo) {
	if len(cst) == 0 {
		return
	}

	g.P("//const as define in tars file")
	g.P("const (")
	for _, v := range g.module.Const {
		v.Rename()
		g.P(v.Name, " ", g.genType(v.Type), " = ", v.Value)
	}
	g.P(")")
}

func (g *GenGo) genInclude(modules []*ast.ModuleInfo) {
	for _, module := range modules {
		genModule := NewGenGo(g.opt, module.Name+module.Source)
		genModule.module = module
		genModule.genAll()
	}
}

func (g *GenGo) genInterface(itf *ast.InterfaceInfo) {
	g.code.Reset()
	itf.Rename()

	g.genHead()
	g.genIFPackage(itf)

	g.genIFServer(itf)
	g.P()
	g.genIFServerWithContext(itf)

	g.genIFProxy(itf)

	g.genIFDispatch(itf)

	g.saveToSourceFile(itf.Name + ".tars.go")
}

func (g *GenGo) genIFProxy(itf *ast.InterfaceInfo) {
	g.P("// ", itf.Name, " struct")
	g.P("type ", itf.Name, ` struct {
	servant model.Servant
}`)

	g.P("// New", itf.Name, " creates a new ", itf.Name, " servant.")
	g.P(`func New`, itf.Name, `() *`, itf.Name, ` {
	return new(`, itf.Name, `)
}`)
	if g.opt.AddServant || !g.opt.WithoutTrace {
		g.P("// New", itf.Name, "Client creates a new ", itf.Name, " client proxy for the given servant.")
		g.P(`func New`, itf.Name, `Client(servant string, comm *tars.Communicator, option ...tars.EndpointManagerOption) *`, itf.Name, ` {
	client := new(`, itf.Name, `)
	comm.StringToProxy(servant, client, option...)
	return client
}`)
	}

	g.P(`// SetServant sets servant for the service.
func (obj *`, itf.Name, `) SetServant(servant model.Servant) {
	obj.servant = servant
}`)

	g.P(`// TarsSetTimeout sets the timeout for the servant which is in ms.
func (obj *`, itf.Name, `) TarsSetTimeout(timeout int) {
	obj.servant.TarsSetTimeout(timeout)
}`)

	g.P(`// TarsSetProtocol sets the protocol for the servant.
func (obj *`, itf.Name, `) TarsSetProtocol(p model.Protocol) {
	obj.servant.TarsSetProtocol(p)
}`)

	g.P(`// Endpoints returns all active endpoint.Endpoint
func (obj *`, itf.Name, `) Endpoints() []*endpoint.Endpoint {
	return obj.servant.Endpoints()
}`)

	if g.opt.AddServant {
		g.P(`// AddServant adds servant  for the service.
func (obj *`, itf.Name, `) AddServant(imp `, itf.Name, `Servant, servant string) {
  tars.AddServant(obj, imp, servant)
}`)
		g.P(`// AddServantWithContext adds servant  for the service with context.
func (obj *`, itf.Name, `) AddServantWithContext(imp `, itf.Name, `ServantWithContext, servant string) {
  tars.AddServantWithContext(obj, imp, servant)
}`)
	}

	for _, v := range itf.Fun {
		g.genIFProxyFun(itf.Name, &v, false, false)
		g.genIFProxyFun(itf.Name, &v, true, false)
		g.genIFProxyFun(itf.Name, &v, true, true)
	}
}

func (g *GenGo) genIFProxyFun(interfName string, fun *ast.FunInfo, withContext bool, isOneWay bool) {
	if withContext {
		if isOneWay {
			g.P("// ", fun.Name, "OneWayWithContext is the proxy function for the method defined in the tars file, with the context")
			g.W("func (obj *", interfName, ") ", fun.Name, "OneWayWithContext(tarsCtx context.Context,")
		} else {
			g.P("// ", fun.Name, "WithContext is the proxy function for the method defined in the tars file, with the context")
			g.W("func (obj *", interfName, ") ", fun.Name, "WithContext(tarsCtx context.Context,")
		}
	} else {
		g.P("// ", fun.Name, " is the proxy function for the method defined in the tars file, with the context")
		g.W("func (obj *", interfName, ") ", fun.Name, "(")
	}
	g.genArgs(fun.Args)

	g.W(" opts ...map[string]string)")

	// not WithContext caller WithContext method
	if !withContext {
		if fun.HasRet {
			g.P("(", g.genType(fun.RetType), ", error) {")
		} else {
			g.P("error {")
		}

		g.W("return obj.", fun.Name, "WithContext(context.Background(), ")
		for _, v := range fun.Args {
			g.W(v.Name, ",")
		}
		g.P(" opts ...)")
		g.P("}")
		return
	}

	if fun.HasRet {
		g.P("(ret ", g.genType(fun.RetType), ", err error) {")
	} else {
		g.P("(err error) {")
	}

	g.P(`var (
		length int32
		have bool
		ty byte
	)`)
	g.P("buf := codec.NewBuffer()")
	var isOut bool
	for k, v := range fun.Args {
		if v.IsOut {
			isOut = true
		}
		dummy := &ast.StructMember{
			Tag:     int32(k + 1),
			Require: true,
			Type:    v.Type,
			Key:     v.Name,
		}
		if v.IsOut {
			dummy.Key = "(*" + dummy.Key + ")"
		}
		g.genWriteVar(dummy, "", fun.HasRet)
		g.P()
	}
	// empty args and below separate
	errStr := errString(fun.HasRet)

	// trace
	if !isOneWay && !g.opt.WithoutTrace {
		g.P(`
trace, ok := current.GetTarsTrace(tarsCtx)
if ok && trace.Call() {
	var traceParam string
	trace.NewSpan()
	traceParamFlag := trace.NeedTraceParam(tarstrace.EstCS, uint(buf.Len()))
	if traceParamFlag == tarstrace.EnpNormal {
		value := map[string]interface{}{}`)
		for _, v := range fun.Args {
			if !v.IsOut {
				g.P("value[", strconv.Quote(v.Name), "] = ", v.Name)
			}
		}
		g.P(`jm, _ := json.Marshal(value)
		traceParam = string(jm)
	} else if traceParamFlag == tarstrace.EnpOverMaxLen {`)
		g.P("traceParam = `{\"trace_param_over_max_len\":true}`")
		g.P(`}
	tars.Trace(trace.GetTraceKey(tarstrace.EstCS), tarstrace.AnnotationCS, tars.GetClientConfig().ModuleName, obj.servant.Name(), `, strconv.Quote(fun.OriginName), `, 0, traceParam, "")
}`)
	}
	g.P()
	g.P(`var statusMap map[string]string
			var contextMap map[string]string
			if len(opts) == 1{
				contextMap =opts[0]
			}else if len(opts) == 2 {
				contextMap = opts[0]
				statusMap = opts[1]
			}
			
			tarsResp := new(requestf.ResponsePacket)`)
	if isOneWay {
		g.P("err = obj.servant.TarsInvoke(tarsCtx, 1, ", strconv.Quote(fun.OriginName), ", buf.ToBytes(), statusMap, contextMap, tarsResp)")
	} else {
		g.P("err = obj.servant.TarsInvoke(tarsCtx, 0, ", strconv.Quote(fun.OriginName), ", buf.ToBytes(), statusMap, contextMap, tarsResp)")
	}
	g.P(errStr)

	if (isOut || fun.HasRet) && !isOneWay {
		g.P("readBuf := codec.NewReader(tools.Int8ToByte(tarsResp.SBuffer))")
	}
	if fun.HasRet && !isOneWay {
		dummy := &ast.StructMember{
			Tag:     0,
			Require: true,
			Type:    fun.RetType,
			Key:     "ret",
		}
		g.genReadVar(dummy, "", fun.HasRet)
		g.P()
	}

	if !isOneWay {
		for k, v := range fun.Args {
			if v.IsOut {
				dummy := &ast.StructMember{
					Tag:     int32(k + 1),
					Require: true,
					Type:    v.Type,
					Key:     "(*" + v.Name + ")",
				}
				g.genReadVar(dummy, "", fun.HasRet)
			}
		}
		if withContext && !g.opt.WithoutTrace {
			traceParamFlag := "traceParamFlag := trace.NeedTraceParam(tarstrace.EstCR, uint(0))"
			if isOut || fun.HasRet {
				traceParamFlag = "traceParamFlag := trace.NeedTraceParam(tarstrace.EstCR, uint(readBuf.Len()))"
			}
			g.P(`
if ok && trace.Call() {
	var traceParam string
	`, traceParamFlag, `
	if traceParamFlag == tarstrace.EnpNormal {
		value := map[string]interface{}{}`)
			if fun.HasRet {
				g.P(`value[""] = ret`)
			}
			for _, v := range fun.Args {
				if v.IsOut {
					g.P("value[", strconv.Quote(v.Name), "] = *", v.Name)
				}
			}
			g.P(`jm, _ := json.Marshal(value)
		traceParam = string(jm)
	} else if traceParamFlag == tarstrace.EnpOverMaxLen {`)
			g.P("traceParam = `{\"trace_param_over_max_len\":true}`")
			g.P(`}
	tars.Trace(trace.GetTraceKey(tarstrace.EstCR), tarstrace.AnnotationCR, tars.GetClientConfig().ModuleName, obj.servant.Name(), `, strconv.Quote(fun.OriginName), `, tarsResp.IRet, traceParam, "")
}`)
			g.P()
		}

		g.P(`
	if len(opts) == 1 {
		for k := range(contextMap){
			delete(contextMap, k)
		}
		for k, v := range(tarsResp.Context){
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range(contextMap){
			delete(contextMap, k)
		}
		for k, v := range(tarsResp.Context){
			contextMap[k] = v
		}
		for k := range(statusMap){
			delete(statusMap, k)
		}
		for k, v := range(tarsResp.Status){
			statusMap[k] = v
		}
	}`)
	}

	g.P()
	g.P(`_ = length
			  _ = have
			  _ = ty`)
	if fun.HasRet {
		g.P("return ret, nil")
	} else {
		g.P("return nil")
	}
	g.P("}")
}

func (g *GenGo) genArgs(args []ast.ArgInfo) {
	for _, arg := range args {
		g.W(arg.Name, " ")
		if arg.IsOut || arg.Type.CType == token.Struct {
			g.W("*")
		}
		g.W(g.genType(arg.Type), ",")
	}
}

func (g *GenGo) genCallArgs(args []ast.ArgInfo) {
	for _, arg := range args {
		if arg.IsOut || arg.Type.CType == token.Struct {
			g.W("&", arg.Name, ",")
		} else {
			g.W(arg.Name, ",")
		}
	}
}

func (g *GenGo) genIFServer(itf *ast.InterfaceInfo) {
	g.P("type ", itf.Name, "Servant interface {")
	for _, v := range itf.Fun {
		g.genIFServerFun(&v)
	}
	g.P("}")
}

func (g *GenGo) genIFServerWithContext(itf *ast.InterfaceInfo) {
	g.P("type ", itf.Name, "ServantWithContext interface {")
	for _, v := range itf.Fun {
		g.genIFServerFunWithContext(&v)
	}
	g.P("}")
}

func (g *GenGo) genIFServerFun(fun *ast.FunInfo) {
	g.W(fun.Name, "(")
	g.genArgs(fun.Args)
	g.W(") (")

	if fun.HasRet {
		g.W("ret ", g.genType(fun.RetType), ", ")
	}
	g.P("err error)")
}

func (g *GenGo) genIFServerFunWithContext(fun *ast.FunInfo) {
	g.W(fun.Name, "(tarsCtx context.Context, ")
	g.genArgs(fun.Args)
	g.W(") (")

	if fun.HasRet {
		g.W("ret ", g.genType(fun.RetType), ", ")
	}
	g.P("err error)")
}

func (g *GenGo) genIFDispatch(itf *ast.InterfaceInfo) {
	g.P("// Dispatch is used to call the server side implement for the method defined in the tars file. withContext shows using context or not.  ")
	g.P("func(obj *", itf.Name, `) Dispatch(tarsCtx context.Context, val interface{}, tarsReq *requestf.RequestPacket, tarsResp *requestf.ResponsePacket, withContext bool) (err error) {
	var (
		length int32
		have bool
		ty byte
	)`)

	var param bool
	for _, v := range itf.Fun {
		if len(v.Args) > 0 {
			param = true
			break
		}
	}

	if param {
		g.P("readBuf := codec.NewReader(tools.Int8ToByte(tarsReq.SBuffer))")
	} else {
		g.P("readBuf := codec.NewReader(nil)")
	}
	g.P(`buf := codec.NewBuffer()
	switch tarsReq.SFuncName {`)
	for _, v := range itf.Fun {
		g.genSwitchCase(itf.Name, &v)
	}

	g.P(`
	default:
		return fmt.Errorf("func mismatch")
	}
	var statusMap map[string]string
	if status, ok := current.GetResponseStatus(tarsCtx); ok  && status != nil {
		statusMap = status
	}
	var contextMap map[string]string
	if ctx, ok := current.GetResponseContext(tarsCtx); ok && ctx != nil  {
		contextMap = ctx
	}
	*tarsResp = requestf.ResponsePacket{
		IVersion:     tarsReq.IVersion,
		CPacketType:  0,
		IRequestId:   tarsReq.IRequestId,
		IMessageType: 0,
		IRet:         0,
		SBuffer:      tools.ByteToInt8(buf.ToBytes()),
		Status:       statusMap,
		SResultDesc:  "",
		Context:      contextMap,
	}

	_ = readBuf
	_ = buf
	_ = length
	_ = have
	_ = ty
	return nil
}`)
}

func (g *GenGo) genSwitchCase(tname string, fun *ast.FunInfo) {
	g.P("case ", strconv.Quote(fun.OriginName), ":")
	inArgsCount := 0
	outArgsCount := 0
	for _, v := range fun.Args {
		g.P("var ", v.Name, " ", g.genType(v.Type))
		if v.Type.Type == token.TMap {
			g.P(v.Name, " = make(", g.genType(v.Type), ")")
		} else if v.Type.Type == token.TVector {
			g.P(v.Name, " = make(", g.genType(v.Type), ", 0)")
		}
		if v.IsOut {
			outArgsCount++
		} else {
			inArgsCount++
		}
	}

	//fmt.Println("args count, in, out:", inArgsCount, outArgsCount)

	if inArgsCount > 0 {
		g.P("if tarsReq.IVersion == basef.TARSVERSION {")
		for k, v := range fun.Args {
			if !v.IsOut {
				dummy := &ast.StructMember{
					Tag:     int32(k + 1),
					Require: true,
					Type:    v.Type,
					Key:     v.Name,
				}
				g.genReadVar(dummy, "", false)
			}
		}

		g.P(`} else if tarsReq.IVersion == basef.TUPVERSION {
		reqTup := tup.NewUniAttribute()
		reqTup.Decode(readBuf)

		var tupBuffer []byte
		`)
		for _, v := range fun.Args {
			if !v.IsOut {
				g.P()
				g.P("reqTup.GetBuffer(", strconv.Quote(v.Name), ", &tupBuffer)")
				g.P("readBuf.Reset(tupBuffer)")

				dummy := &ast.StructMember{
					Tag:     0,
					Require: true,
					Type:    v.Type,
					Key:     v.Name,
				}
				g.genReadVar(dummy, "", false)
			}
		}

		g.P(`} else if tarsReq.IVersion == basef.JSONVERSION {
		var jsonData map[string]interface{}
		decoder := json.NewDecoder(bytes.NewReader(readBuf.ToBytes()))
		decoder.UseNumber()
		err = decoder.Decode(&jsonData)
		if err != nil {
			return fmt.Errorf("decode reqpacket failed, error: %+v", err)
		}`)

		for _, v := range fun.Args {
			if !v.IsOut {
				g.P("{")
				g.P("jsonStr, _ := json.Marshal(jsonData[", strconv.Quote(v.Name), "])")
				if v.Type.CType == token.Struct {
					g.P(v.Name, ".ResetDefault()")
				}
				g.P("if err = json.Unmarshal(jsonStr, &", v.Name, "); err != nil {")
				g.P("return err")
				g.P("}")
				g.P("}")
			}
		}

		g.P(`} else {
			err = fmt.Errorf("decode reqpacket fail, error version: %d", tarsReq.IVersion)
			return err
		}`)

		g.P()
	}
	if !g.opt.WithoutTrace {
		g.P(`
trace, ok := current.GetTarsTrace(tarsCtx)
if ok && trace.Call() {
	var traceParam string
	traceParamFlag := trace.NeedTraceParam(tarstrace.EstSR, uint(readBuf.Len()))
	if traceParamFlag == tarstrace.EnpNormal {
		value := map[string]interface{}{}`)
		for _, v := range fun.Args {
			if !v.IsOut {
				g.P("value[", strconv.Quote(v.Name), "] = ", v.Name)
			}
		}
		g.P(`jm, _ := json.Marshal(value)
		traceParam = string(jm)
	} else if traceParamFlag == tarstrace.EnpOverMaxLen {`)
		g.P("traceParam = `{\"trace_param_over_max_len\":true}`")
		g.P(`}
	tars.Trace(trace.GetTraceKey(tarstrace.EstSR), tarstrace.AnnotationSR, tars.GetClientConfig().ModuleName, tarsReq.SServantName, `, strconv.Quote(fun.OriginName), `, 0, traceParam, "")
}`)
		g.P()
	}

	if fun.HasRet {
		g.P("var funRet ", g.genType(fun.RetType))
		g.P("if !withContext {")
		g.P("imp := val.(", tname, "Servant)")
		g.W("funRet, err = imp.", fun.Name, "(")
		g.genCallArgs(fun.Args)
		g.P(")")

		g.P("} else {")
		g.P("imp := val.(", tname, "ServantWithContext)")
		g.W("funRet, err = imp.", fun.Name, "(tarsCtx ,")
		g.genCallArgs(fun.Args)
		g.P(")")
		g.P("}")
	} else {
		g.P("if !withContext {")
		g.P("imp := val.(", tname, "Servant)")
		g.W("err = imp.", fun.Name, "(")
		g.genCallArgs(fun.Args)
		g.P(")")

		g.P("} else {")
		g.P("imp := val.(", tname, "ServantWithContext)")
		g.W("err = imp.", fun.Name, "(tarsCtx ,")
		g.genCallArgs(fun.Args)
		g.P(")")
		g.P("}")
	}

	if g.opt.DispatchReporter {
		var inArgStr, outArgStr, retArgStr string
		if fun.HasRet {
			retArgStr = "funRet, err"
		} else {
			retArgStr = "err"
		}
		for _, v := range fun.Args {
			prefix := ""
			if v.Type.CType == token.Struct {
				prefix = "&"
			}
			if v.IsOut {
				outArgStr += prefix + v.Name + ","
			} else {
				inArgStr += prefix + v.Name + ","
			}
		}
		g.P("if dp := tars.GetDispatchReporter(); dp != nil {")
		g.P("dp(tarsCtx, []interface{}{", inArgStr, "}, []interface{}{", outArgStr, "}, []interface{}{", retArgStr, "})")
		g.P("}")
	}
	g.P(errString(false))

	g.P()
	g.P("if tarsReq.IVersion == basef.TARSVERSION {")
	g.P("buf.Reset()")
	g.P()
	if fun.HasRet {
		dummy := &ast.StructMember{
			Tag:     0,
			Require: true,
			Type:    fun.RetType,
			Key:     "funRet",
		}
		g.genWriteVar(dummy, "", false)
	}

	for k, v := range fun.Args {
		if v.IsOut {
			dummy := &ast.StructMember{
				Tag:     int32(k + 1),
				Require: true,
				Type:    v.Type,
				Key:     v.Name,
			}
			g.genWriteVar(dummy, "", false)
		}
	}

	g.P("} else if tarsReq.IVersion == basef.TUPVERSION {")
	g.P("rspTup := tup.NewUniAttribute()")
	g.P()
	if fun.HasRet {
		dummy := &ast.StructMember{
			Tag:     0,
			Require: true,
			Type:    fun.RetType,
			Key:     "funRet",
		}
		g.genWriteVar(dummy, "", false)

		g.P(`
		rspTup.PutBuffer("", buf.ToBytes())
		rspTup.PutBuffer("tars_ret", buf.ToBytes())`)
	}

	for _, v := range fun.Args {
		if v.IsOut {
			g.P()
			g.P("buf.Reset()")
			dummy := &ast.StructMember{
				Tag:     0,
				Require: true,
				Type:    v.Type,
				Key:     v.Name,
			}
			g.genWriteVar(dummy, "", false)

			g.P("rspTup.PutBuffer(", strconv.Quote(v.Name), ", buf.ToBytes())")
		}
	}

	g.P(`
	buf.Reset()
	err = rspTup.Encode(buf)
	if err != nil {
		return err
	}
} else if tarsReq.IVersion == basef.JSONVERSION {
	rspJson := map[string]interface{}{}`)
	if fun.HasRet {
		g.P(`rspJson["tars_ret"] = funRet`)
	}

	for _, v := range fun.Args {
		if v.IsOut {
			g.P("rspJson[", strconv.Quote(v.Name), "] = ", v.Name)
		}
	}

	g.P(`
		var rspByte []byte
		if rspByte, err = json.Marshal(rspJson); err != nil {
			return err
		}

		buf.Reset()
		err = buf.WriteSliceUint8(rspByte)
		if err != nil {
			return err
		}
}`)

	if !g.opt.WithoutTrace {
		g.P(`
if ok && trace.Call() {
	var traceParam string
	traceParamFlag := trace.NeedTraceParam(tarstrace.EstSS, uint(buf.Len()))
	if traceParamFlag == tarstrace.EnpNormal {
		value := map[string]interface{}{}`)
		if fun.HasRet {
			g.P(`value[""] = funRet`)
		}
		for _, v := range fun.Args {
			if v.IsOut {
				g.P("value[", strconv.Quote(v.Name), "] = ", v.Name)
			}
		}
		g.P(`jm, _ := json.Marshal(value)
		traceParam = string(jm)
	} else if traceParamFlag == tarstrace.EnpOverMaxLen {`)
		g.P("traceParam = `{\"trace_param_over_max_len\":true}`")
		g.P(`}
	tars.Trace(trace.GetTraceKey(tarstrace.EstSS), tarstrace.AnnotationSS, tars.GetClientConfig().ModuleName, tarsReq.SServantName, `, strconv.Quote(fun.OriginName), `, 0, traceParam, "")
}`)
	}
}
