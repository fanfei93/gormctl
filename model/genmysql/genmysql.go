package genmysql

import (
	"fmt"
	"github.com/xxjwxc/public/mylog"
	"github.com/xxjwxc/public/mysqldb"
	"strings"
)

type MysqlModel struct {
	Orm       *mysqldb.MySqlDB
	TableName string
}

type MysqlOption struct {
	Host      string
	Port      int
	Database  string
	Username  string
	Password  string
	TableName string
}

func NewMysqlModel(option *MysqlOption) *MysqlModel {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&interpolateParams=True",
		option.Username,
		option.Password,
		option.Host,
		option.Port,
		option.Database,
	)
	orm := mysqldb.OnInitDBOrm(dataSource)

	model := &MysqlModel{
		Orm:       orm,
		TableName: option.TableName,
	}

	return model
}

func (m *MysqlModel) Close() {
	m.Orm.OnDestoryDB()
}

func (m *MysqlModel) GetTableBuilderSql() (string, error) {
	rows, err := m.Orm.Raw("show create table " + assemblyTable(m.TableName)).Rows()
	if err != nil {
		fmt.Printf("获取表信息失败，错误：%v\n", err)
		return "", err
	}

	defer rows.Close()

	var tableName, createTable string
	if rows.Next() {
		err := rows.Scan(&tableName, &createTable)
		if err != nil {
			fmt.Printf("获取表信息失败，rows.Scan执行失败，错误：%v\n", err)
			return "", err
		}

		fmt.Printf("tableName: %v, createTable: %v\n", tableName, createTable)
	}

	return createTable, nil
}

func (m *MysqlModel) GetTableElement() (el []ColumnsInfo, err error) {
	keyNameCount := make(map[string]int)
	KeyColumnMp := make(map[string][]keys)
	// get keys
	var Keys []keys
	m.Orm.Raw("show keys from " + assemblyTable(m.TableName)).Scan(&Keys)
	for _, v := range Keys {
		keyNameCount[v.KeyName]++
		KeyColumnMp[v.ColumnName] = append(KeyColumnMp[v.ColumnName], v)
	}
	// ----------end

	var list []genColumns
	// Get table annotations.获取表注释
	m.Orm.Raw("show FULL COLUMNS from " + assemblyTable(m.TableName)).Scan(&list)
	// filter gorm.Model.过滤 gorm.Model
	if filterModel(&list) {
		el = append(el, ColumnsInfo{
			Type: "gorm.Model",
		})
	}
	// -----------------end

	for _, v := range list {
		var tmp ColumnsInfo
		tmp.Name = v.Field
		tmp.Type = v.Type
		tmp.Extra = v.Extra
		FixNotes(&tmp, v.Desc) // 分析表注释

		if v.Default != nil {
			if *v.Default == "" {
				tmp.Gormt = "default:''"
			} else {
				tmp.Gormt = fmt.Sprintf("default:%s", *v.Default)
			}
		}

		tmp.IsNull = strings.EqualFold(v.Null, "YES")
		el = append(el, tmp)
	}
	return
}

func assemblyTable(name string) string {
	return "`" + name + "`"
}

// filterModel filter.过滤 gorm.Model
func filterModel(list *[]genColumns) bool {
	var _temp []genColumns
	num := 0
	for _, v := range *list {
		if strings.EqualFold(v.Field, "id") ||
			strings.EqualFold(v.Field, "created_at") ||
			strings.EqualFold(v.Field, "updated_at") ||
			strings.EqualFold(v.Field, "deleted_at") {
			num++
		} else {
			_temp = append(_temp, v)
		}
	}

	if num >= 4 {
		*list = _temp
		return true
	}

	return false
}

// FixNotes 分析元素表注释
func FixNotes(em *ColumnsInfo, note string) {
	b0 := FixElementTag(em, note) // gorm
	if !b0 {                      // 补偿
		FixElementTag(em, em.Notes) // gorm
	}
}

// FixElementTag 分析元素表注释
func FixElementTag(em *ColumnsInfo, note string) bool {
	matches := noteRegex.FindStringSubmatch(note)
	if len(matches) < 2 {
		em.Notes = note
		return false
	}

	mylog.Infof("get one gorm tag:(%v) ==> (%v)", em.BaseInfo.Name, matches[1])
	em.Notes = note[len(matches[0]):]
	em.Gormt = matches[1]
	return true
}
