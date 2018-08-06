package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
)

const ProjDirInGoPath = "src/github.com/andrskom/gorpcgen"

func main() {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		var err error
		goPath, err = filepath.Abs("~/go")
		if err != nil {
			log.WithError(err).Fatal("Can't get GOPATH")
		}
	}

	var (
		genPath       string
		serverSubPath string
		clientSubPath string
		useValidation bool
		serverTmpl    string
		clientTmpl    string
		handlersPath  string
	)
	flag.StringVar(&genPath, "gen.path", "gen", "Path for generating")
	flag.StringVar(&serverSubPath, "gen.server-path", "server", "Subpath for generating server")
	flag.StringVar(&clientSubPath, "gen.client-path", "client", "Subpath for generating client")
	flag.BoolVar(&useValidation, "cfg.use-validation", false, "Set if you want to use validation")
	flag.StringVar(&handlersPath, "cfg.handlers-path", "./pkg/service/handlers", "Path to handlers")
	flag.StringVar(
		&serverTmpl,
		"cfg.server-tmpl",
		filepath.Join(goPath, ProjDirInGoPath, ".templates/nats/server.gotmpl"),
		"Set path to server tmplt",
	)
	flag.StringVar(
		&clientTmpl,
		"cfg.client-tmpl",
		filepath.Join(goPath, ProjDirInGoPath, ".templates/nats/client.gotmpl"),
		"Set path to client tmplt",
	)
	flag.Parse()

	projectPath, err := os.Getwd()
	if err != nil {
		log.WithError(err).Fatal("Can't get work dir")
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, handlersPath, filter, parser.ParseComments)
	if err != nil {
		log.WithError(err).Fatal("Can't parse package")
	}
	if len(pkgs) != 1 {
		log.Fatal("Expected only one package for parsing")
	}

	parseRPC := &parseRPCStruct{
		HandlersPackage: filepath.Join(projectPath[len(goPath)+5:], handlersPath),
		UseValidation:   useValidation,
		ServerPackage:   strings.Split(serverSubPath, "/")[len(strings.Split(serverSubPath, "/"))-1],
		ClientPackage:   strings.Split(clientSubPath, "/")[len(strings.Split(clientSubPath, "/"))-1],
		Imports:         make(map[string]string),
	}

	parseRPC.regExpService, err = regexp.Compile("//[ ]*goRpcGen:service")
	if err != nil {
		log.WithError(err).Fatal("Can't parse package")
	}

	parseRPC.regExpMethod, err = regexp.Compile("//[ ]*goRpcGen:method")
	if err != nil {
		log.WithError(err).Fatal("Can't parse package")
	}

	for pkgName, astPackage := range pkgs {
		parseRPC.Name = pkgName
		parseRPC.parsePackage(astPackage)
	}

	parseRPC.sortMethods()

	if err := updateServer(filepath.Join(projectPath, genPath, serverSubPath), serverTmpl, parseRPC); err != nil {
		log.WithError(err).Error("Can't generate service server")
	}

	if err := updateClient(filepath.Join(projectPath, genPath, clientSubPath), clientTmpl, parseRPC); err != nil {
		log.WithError(err).Error("Can't generate service client")
	}
}

type parseRPCStruct struct {
	regExpService   *regexp.Regexp
	regExpMethod    *regexp.Regexp
	UseValidation   bool
	HandlersPackage string
	ServerPackage   string
	ClientPackage   string
	Name            string
	ServiceName     string
	Methods         []Method
	Imports         map[string]string
}

func (pr *parseRPCStruct) sortMethods() {
	methods := make(map[string]Method)
	methodsName := make([]string, 0)
	for _, v := range pr.Methods {
		methods[v.Name] = v
		methodsName = append(methodsName, v.Name)
	}
	sort.Strings(methodsName)
	sorted := make([]Method, 0)
	for _, n := range methodsName {
		sorted = append(sorted, methods[n])
	}
	pr.Methods = sorted
}

// Method is info about method for generating
type Method struct {
	Name         string
	TitleName    string
	RequestName  string
	ResponseName string
}

func (pr *parseRPCStruct) parsePackage(p *ast.Package) {
	for _, file := range p.Files {
		pr.parseFile(file)
	}
}

func (pr *parseRPCStruct) parseFile(f *ast.File) {
	for _, decl := range f.Decls {
		if reflect.TypeOf(decl).Elem().String() == "ast.GenDecl" {
			genDecl := decl.(*ast.GenDecl)
			if !(genDecl.Tok == token.TYPE && docMatch(pr.regExpService, genDecl.Doc)) {
				continue
			}
			if !(len(genDecl.Specs) == 1 && reflect.TypeOf(genDecl.Specs[0]).Elem().String() == "ast.TypeSpec") {
				continue
			}
			spec := genDecl.Specs[0].(*ast.TypeSpec)
			pr.ServiceName = spec.Name.Name
		} else if reflect.TypeOf(decl).Elem().String() == "ast.FuncDecl" {
			funcDecl := decl.(*ast.FuncDecl)
			if !docMatch(pr.regExpMethod, funcDecl.Doc) {
				continue
			}
			m := Method{Name: funcDecl.Name.Name, TitleName: strings.Title(funcDecl.Name.Name)}
			if funcDecl.Type != nil &&
				funcDecl.Type.Params != nil &&
				len(funcDecl.Type.Params.List) != 2 &&
				reflect.TypeOf(funcDecl.Type.Params.List[1].Type).String() != "*ast.StarExpr" {

				log.Errorf("bad implementation of rpc method '%s'", m.Name)
				continue
			}
			if funcDecl.Type != nil &&
				funcDecl.Type.Results != nil &&
				len(funcDecl.Type.Results.List) != 2 &&
				reflect.TypeOf(funcDecl.Type.Results.List[0].Type).String() != "*ast.StarExpr" {

				log.Errorf("bad implementation of rpc method '%s'", m.Name)
				continue
			}
			m.RequestName = funcDecl.Type.Params.List[1].Type.(*ast.StarExpr).X.(*ast.Ident).Name
			m.ResponseName = funcDecl.Type.Results.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
			pr.Methods = append(pr.Methods, m)
		}
	}
}

func docMatch(regexp *regexp.Regexp, doc *ast.CommentGroup) bool {
	res := false
	if doc != nil {
		for _, cm := range doc.List {
			if regexp.Match([]byte(cm.Text)) {
				res = true
			}
		}
	}

	return res
}

func updateServer(path string, tmplPath string, parseRPC *parseRPCStruct) error {
	file, err := os.OpenFile(filepath.Join(path, "server.go"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			log.WithError(errClose).Error("Can't close file")
		}
	}()

	tmplString, err := LoadTmpltData(tmplPath)
	tmpl, err := template.New("service").Parse(tmplString)
	if err != nil {
		return err
	}

	return tmpl.Execute(file, parseRPC)
}

func updateClient(path string, tmplPath string, parseRPC *parseRPCStruct) error {
	file, err := os.OpenFile(filepath.Join(path, "client.go"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			log.WithError(errClose).Error("Can't close file")
		}
	}()
	tmplString, err := LoadTmpltData(tmplPath)
	tmpl, err := template.New("client").Parse(tmplString)
	if err != nil {
		return err
	}

	return tmpl.Execute(file, parseRPC)
}

func filter(fileInfo os.FileInfo) bool {
	if fileInfo.IsDir() {
		return false
	}

	return !strings.HasSuffix(fileInfo.Name(), "_gen.go")
}

func LoadTmpltData(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
