package query

import (
	"fmt"
	"os"
	"text/template"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func GenQuery(dsn, path string) {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "",
		FieldWithIndexTag: true,
		// 表字段默认值与模型结构体字段零值不一致的字段, 在插入数据时需要赋值该字段值为零值的, 结构体字段须是指针类型才能成功, 即`FieldCoverable:true`配置下生成的结构体字段.
		// 因为在插入时遇到字段为零值的会被GORM赋予默认值. 如字段`age`表默认值为10, 即使你显式设置为0最后也会被GORM设为10提交.
		// 如果该字段没有上面提到的插入时赋零值的特殊需要, 则字段为非指针类型使用起来会比较方便.
		FieldCoverable: true,
	})
	db, _ := gorm.Open(mysql.Open(dsn))
	g.UseDB(db)
	tableList, err := db.Migrator().GetTables()
	if err != nil {
		panic(fmt.Errorf("get all tables fail: %w", err))
	}

	for _, tableName := range tableList {
		mate := g.GenerateModel(tableName)
		data := Query{StructName: mate.ModelStructName}
		for _, value := range mate.Fields {
			field := Field{Name: value.Name, Type: value.Type, ColumnName: value.ColumnName,
				ColumnComment: value.ColumnComment, MultilineComment: value.MultilineComment,
				Tag: value.Tag, GORMTag: value.GORMTag, CustomGenType: value.CustomGenType, Relation: value.Relation}
			data.Fields = append(data.Fields, field)
		}
		tpl := template.Must(template.New("query").Parse(`
		// Code generated DO NOT EDIT.
		// Code generated DO NOT EDIT.
		// Code generated DO NOT EDIT.
		package repo
		
		import (
			"context"
			"order/internal/data/model"
			"order/internal/data/query"
			"time"
		
			"github.com/zeromicro/go-zero/core/stores/cache"
			"gorm.io/gen"
			"gorm.io/gorm"
			"gorm.io/gorm/clause"
		)
		
		type {{ .StructName }}Repo struct {
			query *query.Query
			cs    cache.Cache
			times {{.TableName}}Time
		}
		
		type {{.TableName}}Time struct {
			{{- range $index, $value := .Fields -}}
			{{ if eq $value.Type "*time.Time" }}
			begin{{$.StructName}}{{$value.Name}} time.Time
			end{{$.StructName}}{{$value.Name}} time.Time
			{{ end }}
			{{- end }}
		}
		
		
		{{- range $index, $value := .Fields -}}
		{{ if eq $value.Type "*time.Time" }}
		func (t {{$.StructName}}Repo) With{{$value.Name}}(begin, end time.Time) {{ $.StructName }}Repo {
			t.times.begin{{$.StructName}}{{$value.Name}} = begin
			t.times.end{{$.StructName}}{{$value.Name}} = end
		
			return t
		} 
		
		{{ end }}
		{{- end -}}
		
		
		func New{{.StructName}}Repo(db *gorm.DB, cs cache.Cache) {{.StructName}}Repo {
			return {{.StructName}}Repo{query: query.Use(db), cs: cs}
		}
		
		func (t {{.StructName}}Repo) FindList(ctx context.Context, m{{.StructName}} model.{{.StructName}}, offset, limit int) ([]*model.{{.StructName}}, int64, error) {
		   return t.find(ctx, m{{.StructName}}, offset, limit)
		}
		
		func (t {{.StructName}}Repo) FindOne(ctx context.Context, m{{.StructName}} model.{{.StructName}}) (*model.{{.StructName}}, error) {
		   items, _, err := t.find(ctx, m{{.StructName}}, 0, 0)
		   return items[0], err
		}
		
		func (t {{.StructName}}Repo) find(ctx context.Context, m{{.StructName}} model.{{.StructName}}, offset, limit int) ([]*model.{{.StructName}}, int64, error) {
			repo := t.query.{{.StructName}}
			dbq := repo.WithContext(ctx)
			{{- range $index, $value := .Fields -}}
			{{ if eq $value.Type "*string" "string" }}
			if m{{$.StructName}}.{{$value.Name}} != "" {
				dbq.Where(repo.{{$value.Name}}.Eq(m{{$.StructName}}.{{$value.Name}}))
			}
			{{ else if eq $value.Type "*time.Time" }}
			if !t.times.begin{{$.StructName}}{{$value.Name}}.IsZero() {
				dbq.Where(repo.{{$value.Name}}.Lte(t.times.begin{{$.StructName}}{{$value.Name}}))
			}
			if !t.times.end{{$.StructName}}{{$value.Name}}.IsZero() {
				dbq.Where(repo.{{$value.Name}}.Gte(t.times.end{{$.StructName}}{{$value.Name}}))
			}
			{{ else }}
			if m{{$.StructName}}.{{$value.Name}} > 0 {
				dbq.Where(repo.{{$value.Name}}.Eq(m{{$.StructName}}.{{$value.Name}}))
			}
			{{ end }}
			{{- end -}}
		
			if offset == 0 && limit == 0 {
				m, err := dbq.Order(repo.ID.Desc()).First()
				return []*model.{{.StructName}}{m}, 0, err
			} else {
				return dbq.Order(repo.ID.Desc()).FindByPage(offset, limit)
			}
		}
		
		func (t {{.StructName}}Repo) FindById(ctx context.Context, id int64) (*model.{{.StructName}}, error) {
			repo := t.query.{{.StructName}}
			return repo.WithContext(ctx).Where(repo.ID.Eq(id)).Take()
		}
		
		func (t {{.StructName}}Repo) Create(ctx context.Context, m ...*model.{{.StructName}}) error {
			return t.query.{{.StructName}}.WithContext(ctx).Create(m...)
		}
		
		func (t {{.StructName}}Repo) Del(ctx context.Context, id int64) error {
			repo := t.query.{{.StructName}}
			_, err := repo.WithContext(ctx).Where(repo.ID.Eq(id)).Delete()
			return err
		}
		
		func (t {{.StructName}}Repo) Update(ctx context.Context, m{{.StructName}} model.{{.StructName}}) (gen.ResultInfo, error) {
			return t.query.{{.StructName}}.WithContext(ctx).Updates(m{{.StructName}})
		}
		
		func (t {{.StructName}}Repo) CreateOrUpdate(ctx context.Context, m{{.StructName}} model.{{.StructName}}) error {
			datetime := time.Now()
			m{{.StructName}}.CreateTime = datetime
			m{{.StructName}}.UpdateTime = datetime
			data := map[string]any{
				{{- range $index, $value := .Fields -}}
				"{{$value.ColumnName}}": m{{$.StructName}}.{{$value.Name}},
				{{- end -}}
			}
			pk := clause.Column{Name: "id"}
			return t.query.{{.StructName}}.WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{pk},
				DoUpdates: clause.Assignments(data), // 更新哪些字段
			}).Create(&m{{.StructName}})
		}
		

		`))
		if err != nil {
			panic(err)
		}
		filePath := fmt.Sprintf("%s/%s.gen.go", path, mate.TableName)
		file, err := os.OpenFile(filePath, os.O_EXCL|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("文件打开失败", err)
			continue
		}
		//及时关闭file句柄
		defer file.Close()
		err = tpl.Execute(file, data)
		if err != nil {
			panic(err)
		}
	}
}

type Query struct {
	StructName string
	Fields     []Field
}

type Field struct {
	Name             string
	Type             string
	ColumnName       string
	ColumnComment    string
	MultilineComment bool
	Tag              field.Tag
	GORMTag          field.GormTag
	CustomGenType    string
	Relation         *field.Relation
}
