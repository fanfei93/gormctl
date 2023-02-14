package model

import (
	"bytes"
	"fmt"
	"github.com/xxjwxc/public/mybigcamel"
	"gitlab.2345.cn/gomod/gormctl/config"
	"gitlab.2345.cn/gomod/gormctl/model/genmysql"
	gentemplate "gitlab.2345.cn/gomod/gormctl/model/template"
	"gitlab.2345.cn/gomod/gormctl/view/genstruct"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
)

type GenModel struct {
	genmysql.MysqlOption

	PackageName string
}

type GenModelOption struct {
	genmysql.MysqlOption

	PackageName string
}

func NewGenModel(option *GenModelOption) *GenModel {
	return &GenModel{
		MysqlOption: option.MysqlOption,
		PackageName: option.PackageName,
	}
}

func (m *GenModel) GetGenPackage() *genstruct.GenPackage {
	pkg := new(genstruct.GenPackage)
	pkg.Name = m.PackageName

	return pkg
}

func (m *GenModel) getTableInfoStruct() (*gentemplate.GenBaseStruct, error) {
	model := genmysql.NewMysqlModel(&m.MysqlOption)
	elements, err := model.GetTableElement()
	if err != nil {
		return nil, err
	}

	var columns []gentemplate.Column
	var gormModel bool
	for _, element := range elements {
		var columnTag, columnComment string
		if element.Type == "gorm.Model" {
			gormModel = true
		} else {
			columnTag = "`gorm:\"column:" + element.Name + "\"" + " json:\"" + element.Name + "\"`"
			columnComment = "// " + element.Notes
		}

		columnInfo := gentemplate.Column{
			ColumnName: getCamelName(element.Name),
			ColumnType: getTypeName(element.Type, element.IsNull),
			ColumnTag:  columnTag,
			Comment:    columnComment,
		}
		columns = append(columns, columnInfo)
	}

	structName := getCamelName(m.TableName)
	outerInterfaceName := structName + "Model"
	innerInterfaceName := strings.ToLower(outerInterfaceName[:1]) + outerInterfaceName[1:]
	defaultModelName := "default" + outerInterfaceName

	data := &gentemplate.GenBaseStruct{
		PackageName:        m.PackageName,
		InnerInterfaceName: innerInterfaceName,
		OuterInterfaceName: outerInterfaceName,
		DefaultModelName:   defaultModelName,
		StructName:         structName,
		Columns:            columns,
		TableName:          m.TableName,
		GormModel:          gormModel,
	}
	return data, nil
}

func (m *GenModel) GenerateModelFile(c config.Config) error {
	infoStruct, err := m.getTableInfoStruct()
	if err != nil {
		return err
	}
	content, err := m.getBaseModelContent(infoStruct)
	if err != nil {
		return err
	}

	if !isExist(c.OutDir) {
		if !createDir(c.OutDir) {
			os.Exit(1)
		}
	}
	baseModelPath := c.OutDir + "/" + infoStruct.TableName + "_model_gen.go"
	err = ioutil.WriteFile(baseModelPath, content, 0666)
	if err != nil {
		return err
	}
	exec.Command("goimports", "-l", "-w", baseModelPath).Output()
	exec.Command("gofmt", "-l", "-w", baseModelPath).Output()

	fmt.Println(baseModelPath + " Done")

	content, err = m.getCustomModelContent(infoStruct)
	if err != nil {
		return err
	}
	customModelPath := c.OutDir + "/" + infoStruct.TableName + "_model.go"
	if isExist(customModelPath) {
		return nil
	}
	exec.Command("goimports", "-l", "-w", customModelPath).Output()
	exec.Command("gofmt", "-l", "-w", customModelPath).Output()

	err = ioutil.WriteFile(customModelPath, content, 0666)
	if err != nil {
		return err
	}

	fmt.Println(customModelPath + " Done")
	return nil
}

func (m *GenModel) getBaseModelContent(data *gentemplate.GenBaseStruct) ([]byte, error) {
	parse, err := template.New("gen_base").Parse(gentemplate.GetGenBaseTemplate())
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = parse.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *GenModel) getCustomModelContent(data *gentemplate.GenBaseStruct) ([]byte, error) {
	parse, err := template.New("gen_custom").Parse(gentemplate.GetGenCustomTemplate())
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = parse.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// getTypeName Type acquisition filtering.类型获取过滤
func getTypeName(name string, isNull bool) string {
	// todo 自定义匹配类型

	if name == "gorm.Model" {
		return name
	}

	// Precise matching first.先精确匹配
	if v, ok := genmysql.TypeMysqlDicMp[name]; ok {
		if (name == "timestamp" || name == "datetime") && isNull {
			return "*" + v
		}
		return v
	}

	// Fuzzy Regular Matching.模糊正则匹配
	for _, l := range genmysql.TypeMysqlMatchList {
		if ok, _ := regexp.MatchString(l.Key, name); ok {
			if l.Value == "time.Time" && isNull {
				return "*time.Time"
			}
			return l.Value
		}
	}

	panic(fmt.Sprintf("type (%v) not match in any way.maybe need to add on (https://github.com/xxjwxc/gormt/blob/master/data/view/cnf/def.go)", name))
}

// getCamelName Big Hump or Capital Letter.大驼峰或者首字母大写
func getCamelName(name string) string {
	return mybigcamel.Marshal(strings.ToLower(name))
}
func isExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func createDir(dirName string) bool {
	err := os.MkdirAll(dirName, 0777)
	if err != nil {
		fmt.Printf("%v\n", err)
		return false
	}
	return true
}
