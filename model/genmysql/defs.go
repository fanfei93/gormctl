package genmysql

import "regexp"

// BaseInfo base common attribute. 基础属性
type BaseInfo struct {
	Name  string // table name.表名
	Notes string // table comment . 表注释
}

// ColumnsInfo Columns list .表列信息
type ColumnsInfo struct {
	BaseInfo
	IsNull         bool         // null if db is set null
	Extra          string       // Extra (AUTO_INCREMENT 自增加)
	Type           string       // Type.类型标记
	Gormt          string       // 默认值
	// Index          []KList      // index list.index列表
	// ForeignKeyList []ForeignKey // Foreign key list . 表的外键信息
}

type keys struct {
	NonUnique  int    `gorm:"column:Non_unique"`
	KeyName    string `gorm:"column:Key_name"`
	ColumnName string `gorm:"column:Column_name"`
	IndexType  string `gorm:"column:Index_type"`
}

// genColumns show full columns
type genColumns struct {
	Field   string  `gorm:"column:Field"`
	Type    string  `gorm:"column:Type"`
	Key     string  `gorm:"column:Key"`
	Desc    string  `gorm:"column:Comment"`
	Null    string  `gorm:"column:Null"`
	Extra   string  `gorm:"Extra"`
	Default *string `gorm:"column:Default"`
}

var noteRegex = regexp.MustCompile(`^\[@gorm\s(\S+)+\]`)

// TypeMysqlDicMp Accurate matching type.精确匹配类型
var TypeMysqlDicMp = map[string]string{
	"smallint":            "int16",
	"smallint unsigned":   "uint16",
	"int":                 "int",
	"int unsigned":        "uint",
	"bigint":              "int64",
	"bigint unsigned":     "uint64",
	"mediumint":           "int32",
	"mediumint unsigned":  "uint32",
	"varchar":             "string",
	"char":                "string",
	"date":                "datatypes.Date",
	"datetime":            "time.Time",
	"bit(1)":              "[]uint8",
	"tinyint":             "int8",
	"tinyint unsigned":    "uint8",
	"tinyint(1)":          "bool", // tinyint(1) 默认设置成bool
	"tinyint(1) unsigned": "bool", // tinyint(1) 默认设置成bool
	"json":                "datatypes.JSON",
	"text":                "string",
	"timestamp":           "time.Time",
	"double":              "float64",
	"double unsigned":     "float64",
	"mediumtext":          "string",
	"longtext":            "string",
	"float":               "float32",
	"float unsigned":      "float32",
	"tinytext":            "string",
	"enum":                "string",
	"time":                "time.Time",
	"tinyblob":            "[]byte",
	"blob":                "[]byte",
	"mediumblob":          "[]byte",
	"longblob":            "[]byte",
	"integer":             "int64",
	"numeric":             "float64",
	"smalldatetime":       "time.Time", //sqlserver
	"nvarchar":            "string",
	"real":                "float32",
	"binary":              "[]byte",
}

// TypeMysqlMatchList Fuzzy Matching Types.模糊匹配类型
var TypeMysqlMatchList = []struct {
	Key   string
	Value string
}{
	{`^(tinyint)[(]\d+[)] unsigned`, "uint8"},
	{`^(smallint)[(]\d+[)] unsigned`, "uint16"},
	{`^(int)[(]\d+[)] unsigned`, "uint32"},
	{`^(bigint)[(]\d+[)] unsigned`, "uint64"},
	{`^(float)[(]\d+,\d+[)] unsigned`, "float64"},
	{`^(double)[(]\d+,\d+[)] unsigned`, "float64"},
	{`^(tinyint)[(]\d+[)]`, "int8"},
	{`^(smallint)[(]\d+[)]`, "int16"},
	{`^(int)[(]\d+[)]`, "int"},
	{`^(bigint)[(]\d+[)]`, "int64"},
	{`^(char)[(]\d+[)]`, "string"},
	{`^(enum)[(](.)+[)]`, "string"},
	{`^(varchar)[(]\d+[)]`, "string"},
	{`^(varbinary)[(]\d+[)]`, "[]byte"},
	{`^(blob)[(]\d+[)]`, "[]byte"},
	{`^(binary)[(]\d+[)]`, "[]byte"},
	{`^(decimal)[(]\d+,\d+[)]`, "float64"},
	{`^(mediumint)[(]\d+[)]`, "int16"},
	{`^(mediumint)[(]\d+[)] unsigned`, "uint16"},
	{`^(double)[(]\d+,\d+[)]`, "float64"},
	{`^(float)[(]\d+,\d+[)]`, "float64"},
	{`^(datetime)[(]\d+[)]`, "time.Time"},
	{`^(bit)[(]\d+[)]`, "[]uint8"},
	{`^(text)[(]\d+[)]`, "string"},
	{`^(integer)[(]\d+[)]`, "int"},
	{`^(timestamp)[(]\d+[)]`, "time.Time"},
	{`^(geometry)[(]\d+[)]`, "[]byte"},
	{`^(set)[(][\s\S]+[)]`, "string"},
}