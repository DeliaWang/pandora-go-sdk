package pipeline

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"

	"github.com/qiniu/pandora-go-sdk/base"
	"github.com/qiniu/pandora-go-sdk/base/reqerr"
)

type PipelineToken struct {
	Token string `json:"-"`
}

const (
	schemaKeyPattern      = "^[a-zA-Z_][a-zA-Z0-9_]{0,127}$"
	groupNamePattern      = "^[a-zA-Z_][a-zA-Z0-9_]{0,127}$"
	repoNamePattern       = "^[a-zA-Z_][a-zA-Z0-9_]{0,127}$"
	transformNamePattern  = "^[a-zA-Z_][a-zA-Z0-9_]{0,127}$"
	exportNamePattern     = "^[a-zA-Z_][a-zA-Z0-9_]{0,127}$"
	datasourceNamePattern = "^[a-zA-Z_][a-zA-Z0-9_]{0,127}$"
	jobNamePattern        = "^[a-zA-Z_][a-zA-Z0-9_]{0,127}$"
	jobExportNamePattern  = "^[a-zA-Z_][a-zA-Z0-9_]{0,127}$"
	pluginNamePattern     = "^[a-zA-Z][a-zA-Z0-9_\\.]{0,127}[a-zA-Z0-9_]$"
)

var schemaTypes = map[string]bool{
	"float":   true,
	"string":  true,
	"long":    true,
	"date":    true,
	"array":   true,
	"map":     true,
	"boolean": true,
}

func validateGroupName(g string) error {
	matched, err := regexp.MatchString(groupNamePattern, g)
	if err != nil {
		return reqerr.NewInvalidArgs("GroupName", err.Error())
	}
	if !matched {
		return reqerr.NewInvalidArgs("GroupName", fmt.Sprintf("invalid group name: %s", g))
	}
	return nil
}

func validateRepoName(r string) error {
	matched, err := regexp.MatchString(repoNamePattern, r)
	if err != nil {
		return reqerr.NewInvalidArgs("RepoName", err.Error())
	}
	if !matched {
		return reqerr.NewInvalidArgs("RepoName", fmt.Sprintf("invalid repo name: %s", r))
	}
	return nil
}

func validateTransformName(t string) error {
	matched, err := regexp.MatchString(transformNamePattern, t)
	if err != nil {
		return reqerr.NewInvalidArgs("TransformName", err.Error())
	}
	if !matched {
		return reqerr.NewInvalidArgs("TransformName", fmt.Sprintf("invalid transform name: %s", t))
	}
	return nil
}

func validateExportName(e string) error {
	matched, err := regexp.MatchString(exportNamePattern, e)
	if err != nil {
		return reqerr.NewInvalidArgs("ExportName", err.Error())
	}
	if !matched {
		return reqerr.NewInvalidArgs("ExportName", fmt.Sprintf("invalid export name: %s", e))
	}
	return nil
}

func validatePluginName(p string) error {
	matched, err := regexp.MatchString(pluginNamePattern, p)
	if err != nil {
		return reqerr.NewInvalidArgs("PluginName", err.Error())
	}
	if !matched {
		return reqerr.NewInvalidArgs("PluginName", fmt.Sprintf("invalid plugin name: %s", p))
	}
	return nil
}

func validateDatasouceName(d string) error {
	matched, err := regexp.MatchString(datasourceNamePattern, d)
	if err != nil {
		return reqerr.NewInvalidArgs("DatasourceName", err.Error())
	}
	if !matched {
		return reqerr.NewInvalidArgs("DatasourceName", fmt.Sprintf("invalid datasource name: %s", d))
	}
	return nil
}

func validateJobName(j string) error {
	matched, err := regexp.MatchString(datasourceNamePattern, j)
	if err != nil {
		return reqerr.NewInvalidArgs("JobName", err.Error())
	}
	if !matched {
		return reqerr.NewInvalidArgs("JobName", fmt.Sprintf("invalid job name: %s", j))
	}
	return nil
}

func validateJobexportName(e string) error {
	matched, err := regexp.MatchString(datasourceNamePattern, e)
	if err != nil {
		return reqerr.NewInvalidArgs("JobexportName", err.Error())
	}
	if !matched {
		return reqerr.NewInvalidArgs("JobexportName", fmt.Sprintf("invalid job export name: %s", e))
	}
	return nil
}

type Container struct {
	Type   string `json:"type"`
	Count  int    `json:"count"`
	Status string `json:"status,omitempty"`
}

func (c *Container) Validate() (err error) {
	if c.Type != "M16C4" && c.Type != "M32C8" {
		err = reqerr.NewInvalidArgs("ContainerType", fmt.Sprintf("invalid container type: %s, should be one of \"M16C4\" and \"M32C8\"", c.Type))
		return
	}
	if c.Count < 1 || c.Count > 128 {
		err = reqerr.NewInvalidArgs("ContainerCount", fmt.Sprintf("invalid container count: %d", c.Count))
		return
	}
	return
}

type CreateGroupInput struct {
	PipelineToken
	GroupName       string     `json:"-"`
	Region          string     `json:"region"`
	Container       *Container `json:"container"`
	AllocateOnStart bool       `json:"allocateOnStart,omitempty"`
}

func (g *CreateGroupInput) Validate() (err error) {
	if err = validateGroupName(g.GroupName); err != nil {
		return
	}
	if g.Region == "" {
		err = reqerr.NewInvalidArgs("Region", "region should not be empty")
		return
	}
	if g.Container == nil {
		err = reqerr.NewInvalidArgs("Container", "container should not be empty")
		return
	}
	if err = g.Container.Validate(); err != nil {
		return
	}
	return
}

type UpdateGroupInput struct {
	PipelineToken
	GroupName string     `json:"-"`
	Container *Container `json:"container"`
}

func (g *UpdateGroupInput) Validate() (err error) {
	if err = validateGroupName(g.GroupName); err != nil {
		return
	}
	if g.Container == nil {
		err = reqerr.NewInvalidArgs("Container", "container should not be empty")
		return
	}
	if err = g.Container.Validate(); err != nil {
		return
	}
	return
}

type StartGroupTaskInput struct {
	PipelineToken
	GroupName string
}

type StopGroupTaskInput struct {
	PipelineToken
	GroupName string
}

type GetGroupInput struct {
	PipelineToken
	GroupName string
}

type GetGroupOutput struct {
	Region     string     `json:"region"`
	Container  *Container `json:"container"`
	CreateTime string     `json:"createTime"`
	UpdateTime string     `json:"updateTime"`
}

type DeleteGroupInput struct {
	PipelineToken
	GroupName string
}

type GroupDesc struct {
	GroupName string     `json:"name"`
	Region    string     `json:"region"`
	Container *Container `json:"container"`
}

type ListGroupsInput struct {
	PipelineToken
}

type ListGroupsOutput struct {
	Groups []GroupDesc `json:"groups"`
}

type RepoSchemaEntry struct {
	Key       string            `json:"key"`
	ValueType string            `json:"valtype"`
	Required  bool              `json:"required"`
	ElemType  string            `json:"elemtype,omitempty"`
	Schema    []RepoSchemaEntry `json:"schema,omitempty"`
}

func (e RepoSchemaEntry) String() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}

func (e *RepoSchemaEntry) Validate() (err error) {
	matched, err := regexp.MatchString(schemaKeyPattern, e.Key)
	if err != nil {
		err = reqerr.NewInvalidArgs("Schema", err.Error())
		return
	}
	if !matched {
		err = reqerr.NewInvalidArgs("Schema", fmt.Sprintf("invalid field key: %s", e.Key))
		return

	}
	if !schemaTypes[e.ValueType] {
		err = reqerr.NewInvalidArgs("Schema", fmt.Sprintf("invalid field type: %s, field type should be one of \"float\", \"string\", \"date\", \"long\", \"boolean\", \"array\" and \"map\"", e.ValueType))
		return
	}
	if e.ValueType == "array" {
		if e.ElemType != "float" && e.ElemType != "long" && e.ElemType != "string" {
			err = reqerr.NewInvalidArgs("Schema", fmt.Sprintf("invalid field type in array: %s, field type should be one of \"float\", \"string\", and \"long\"", e.ValueType))
			return
		}
	}
	if e.ValueType == "map" {
		for _, ns := range e.Schema {
			if err = ns.Validate(); err != nil {
				return
			}
		}
	}

	return
}

type CreateRepoDSLInput struct {
	PipelineToken
	RepoName  string
	Region    string `json:"region"`
	DSL       string `json:"dsl"`
	GroupName string `json:"group"`
}

/*
DSL创建的规则为`<字段名称> <类型>`,字段名称和类型用空格符隔开，不同字段用逗号隔开。若字段必填，则在类型前加`*`号表示。
    * pandora date类型：`date`,`DATE`,`d`,`D`
    * pandora long类型：`long`,`LONG`,`l`,`L`
    * pandora float类型: `float`,`FLOAT`,`F`,`f`
    * pandora string类型: `string`,`STRING`,`S`,`s`
    * pandora bool类型:  `bool`,`BOOL`,`B`,`b`,`boolean`
    * pandora array类型: `array`,`ARRAY`,`A`,`a`;括号中跟具体array元素的类型，如a(l)，表示array里面都是long。
    * pandora map类型: `map`,`MAP`,`M`,`m`;使用花括号表示具体类型，表达map里面的元素，如map{a l,b map{c b,x s}}, 表示map结构体里包含a字段，类型是long，b字段又是一个map，里面包含c字段，类型是bool，还包含x字段，类型是string。
*/

func getRawType(tp string) (schemaType string, err error) {
	schemaType = strings.ToLower(tp)
	switch schemaType {
	case "l", "long":
		schemaType = "long"
	case "f", "float":
		schemaType = "float"
	case "s", "string":
		schemaType = "string"
	case "d", "date":
		schemaType = "date"
	case "a", "array":
		err = errors.New("arrary type must specify data type surrounded by ( )")
		return
	case "m", "map":
		schemaType = "map"
	case "b", "bool", "boolean":
		schemaType = "boolean"
	case "": //这个是一种缺省
	default:
		err = fmt.Errorf("schema type %v not supperted", schemaType)
		return
	}
	return
}

func getField(f string) (key, valueType, elementType string, required bool, err error) {
	f = strings.TrimSpace(f)
	if f == "" {
		return
	}
	splits := strings.Fields(f)
	switch len(splits) {
	case 1:
		key = splits[0]
		return
	case 2:
		key, valueType = splits[0], splits[1]
	default:
		err = fmt.Errorf("Raw field schema parse error: <%v> was invalid", f)
		return
	}
	if key == "" {
		err = fmt.Errorf("field schema %v key can not be empty", f)
		return
	}
	required = false
	if strings.HasPrefix(valueType, "*") || strings.HasSuffix(valueType, "*") {
		required = true
		valueType = strings.Trim(valueType, "*")
	}
	//处理arrary类型
	if beg := strings.Index(valueType, "("); beg != -1 {
		ed := strings.Index(valueType, ")")
		if ed <= beg {
			err = fmt.Errorf("field schema %v has no type specified", f)
			return
		}
		elementType, err = getRawType(valueType[beg+1 : ed])
		if err != nil {
			err = fmt.Errorf("array 【%v】: %v, key %v valuetype %v", f, err, key, valueType)
		}
		valueType = "array"
		return
	}
	valueType, err = getRawType(valueType)
	if err != nil {
		err = fmt.Errorf("normal 【%v】: %v, key %v valuetype %v", f, err, key, valueType)
	}
	return
}

func toSchema(dsl string, depth int) (schemas []RepoSchemaEntry, err error) {
	if depth > base.NestLimit {
		err = reqerr.NewInvalidArgs("Schema", fmt.Sprintf("RepoSchemaEntry are nested out of limit %v", base.NestLimit))
		return
	}
	schemas = make([]RepoSchemaEntry, 0)
	dsl = strings.TrimSpace(dsl)
	start := 0
	nestbalance := 0
	neststart, nestend := -1, -1
	dsl += "," //增加一个','保证一定是以","为终结
	for end, c := range dsl {
		if start > end {
			err = errors.New("parse dsl inner error: start index is larger than end")
			return
		}
		switch c {
		case '{':
			if nestbalance == 0 {
				neststart = end
			}
			nestbalance++
		case '}':
			nestbalance--
			if nestbalance == 0 {
				nestend = end
				if nestend <= neststart {
					err = errors.New("parse dsl error: nestend should never less or equal than neststart")
					return
				}
				subschemas, err := toSchema(dsl[neststart+1:nestend], depth+1)
				if err != nil {
					return nil, err
				}
				if neststart <= start {
					return nil, errors.New("parse dsl error: map{} not specified")
				}
				key, valueType, _, required, err := getField(dsl[start:neststart])
				if err != nil {
					return nil, err
				}
				if key != "" {
					if valueType == "" {
						valueType = "map"
					}
					schemas = append(schemas, RepoSchemaEntry{
						Key:       key,
						ValueType: valueType,
						Required:  required,
						Schema:    subschemas,
					})
				}
				start = end + 1
			}
		case ',':
			if nestbalance == 0 {
				if start < end {
					key, valueType, elemtype, required, err := getField(strings.TrimSpace(dsl[start:end]))
					if err != nil {
						return nil, err
					}
					if key != "" {
						if valueType == "" {
							valueType = "string"
						}
						schemas = append(schemas, RepoSchemaEntry{
							Key:       key,
							ValueType: valueType,
							Required:  required,
							ElemType:  elemtype,
						})
					}
				}
				start = end + 1
			}
		}
	}
	if nestbalance != 0 {
		err = errors.New("parse dsl error: { and } not match")
		return
	}
	return
}

type CreateRepoInput struct {
	PipelineToken
	RepoName  string
	Region    string            `json:"region"`
	Schema    []RepoSchemaEntry `json:"schema"`
	GroupName string            `json:"group"`
}

func (r *CreateRepoInput) Validate() (err error) {
	if err = validateRepoName(r.RepoName); err != nil {
		return
	}

	if r.Schema == nil || len(r.Schema) == 0 {
		err = reqerr.NewInvalidArgs("Schema", "schema should not be empty")
		return
	}
	for _, schema := range r.Schema {
		if err = schema.Validate(); err != nil {
			return
		}
	}

	if r.GroupName != "" {
		if err = validateGroupName(r.GroupName); err != nil {
			return
		}
	}

	if r.Region == "" {
		err = reqerr.NewInvalidArgs("Region", "region should not be empty")
		return
	}
	return
}

type UpdateRepoInput struct {
	PipelineToken
	RepoName string
	Schema   []RepoSchemaEntry `json:"schema"`
}

func (r *UpdateRepoInput) Validate() (err error) {
	if err = validateRepoName(r.RepoName); err != nil {
		return
	}

	if r.Schema == nil || len(r.Schema) == 0 {
		err = reqerr.NewInvalidArgs("Schema", "schema should not be empty")
		return
	}
	for _, schema := range r.Schema {
		if err = schema.Validate(); err != nil {
			return
		}
	}

	return
}

type GetRepoInput struct {
	PipelineToken
	RepoName string
}

type GetRepoOutput struct {
	Region      string            `json:"region"`
	Schema      []RepoSchemaEntry `json:"schema"`
	GroupName   string            `json:"group"`
	DerivedFrom string            `json:"derivedFrom"`
}

type RepoDesc struct {
	RepoName    string `json:"name"`
	Region      string `json:"region"`
	GroupName   string `json:"group"`
	DerivedFrom string `json:"derivedFrom"`
}

type ListReposInput struct {
	PipelineToken
}

type ListReposOutput struct {
	Repos []RepoDesc `json:"repos"`
}

type DeleteRepoInput struct {
	PipelineToken
	RepoName string
}

type PointField struct {
	Key   string
	Value interface{}
}

func (p *PointField) String() string {
	typ := reflect.TypeOf(p.Value).Kind()
	var value string
	if typ == reflect.Map || typ == reflect.Slice {
		v, _ := json.Marshal(p.Value)
		value = escapeStringField(string(v))
	} else {
		value = escapeStringField(fmt.Sprintf("%v", p.Value))
	}
	return fmt.Sprintf("%s=%s\t", p.Key, value)
}

type Point struct {
	Fields []PointField
}

type Points []Point

func (ps Points) Buffer() []byte {
	var buf bytes.Buffer
	for _, p := range ps {
		for _, field := range p.Fields {
			buf.WriteString(field.String())
		}
		if len(p.Fields) > 0 {
			buf.Truncate(buf.Len() - 1)
		}
		buf.WriteByte('\n')
	}
	if len(ps) > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	return buf.Bytes()
}

func escapeStringField(in string) string {
	var out []byte
	for i := 0; i < len(in); i++ {
		switch in[i] {
		case '\t': // escape tab
			out = append(out, '\\')
			out = append(out, 't')
		case '\n': // escape new line
			out = append(out, '\\')
			out = append(out, 'n')
		default:
			out = append(out, in[i])
		}
	}
	return string(out)
}

type PostDataInput struct {
	PipelineToken
	RepoName string
	Points   Points
}

type PostDataFromFileInput struct {
	PipelineToken
	RepoName string
	FilePath string
}

type PostDataFromReaderInput struct {
	PipelineToken
	RepoName   string
	Reader     io.ReadSeeker
	BodyLength int64
}

type PostDataFromBytesInput struct {
	PipelineToken
	RepoName string
	Buffer   []byte
}

type UploadPluginInput struct {
	PipelineToken
	PluginName string
	Buffer     *bytes.Buffer
}

type UploadPluginFromFileInput struct {
	PipelineToken
	PluginName string
	FilePath   string
}

type GetPluginInput struct {
	PipelineToken
	PluginName string
}

type PluginDesc struct {
	PluginName string `json:"name"`
	CreateTime string `json:"createTime"`
}

type GetPluginOutput struct {
	PluginDesc
}

type ListPluginsInput struct {
	PipelineToken
}

type ListPluginsOutput struct {
	Plugins []PluginDesc `json:"plugins"`
}

type DeletePluginInput struct {
	PipelineToken
	PluginName string
}

type TransformPluginOutputEntry struct {
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
}

type TransformPlugin struct {
	Name   string                       `json:"name"`
	Output []TransformPluginOutputEntry `json:"output"`
}

type TransformSpec struct {
	Plugin    *TransformPlugin `json:"plugin,omitempty"`
	Mode      string           `json:"mode,omitempty"`
	Code      string           `json:"code,omitempty"`
	Interval  string           `json:"interval,omitempty"`
	Container *Container       `json:"container,omitempty"`
}

func (t *TransformSpec) Validate() (err error) {
	if t.Mode == "" && t.Code == "" && t.Plugin == nil {
		err = reqerr.NewInvalidArgs("TransformSpec", "all mode, code and plugin can not be empty")
		return
	}
	if t.Container != nil {
		if err = t.Container.Validate(); err != nil {
			return
		}
	}
	return
}

type CreateTransformInput struct {
	PipelineToken
	SrcRepoName   string
	TransformName string
	DestRepoName  string
	Spec          *TransformSpec
}

func (t *CreateTransformInput) Validate() (err error) {
	if err = validateRepoName(t.SrcRepoName); err != nil {
		return
	}
	if err = validateRepoName(t.DestRepoName); err != nil {
		return
	}
	if err = validateTransformName(t.TransformName); err != nil {
		return
	}
	if t.SrcRepoName == t.DestRepoName {
		err = reqerr.NewInvalidArgs("DestRepoName", "dest repo name should be different to src repo name")
		return
	}
	return t.Spec.Validate()
}

type UpdateTransformInput struct {
	PipelineToken
	SrcRepoName   string
	TransformName string
	Spec          *TransformSpec
}

func (t *UpdateTransformInput) Validate() (err error) {
	if err = validateRepoName(t.SrcRepoName); err != nil {
		return
	}
	if err = validateTransformName(t.TransformName); err != nil {
		return
	}
	return t.Spec.Validate()
}

type TransformDesc struct {
	TransformName string         `json:"name"`
	DestRepoName  string         `json:"to"`
	Spec          *TransformSpec `json:"spec"`
}

type GetTransformInput struct {
	PipelineToken
	RepoName      string
	TransformName string
}

type GetTransformOutput struct {
	TransformDesc
}

type DeleteTransformInput struct {
	PipelineToken
	RepoName      string
	TransformName string
}

type ListTransformsInput struct {
	PipelineToken
	RepoName string
}

type ListTransformsOutput struct {
	Transforms []TransformDesc `json:"transforms"`
}

type ExportFilter struct {
	Rules     map[string]map[string]string `json:"rules"`
	ToDefault bool                         `json:"toDefault"`
}

func (f *ExportFilter) Validate() (err error) {
	if len(f.Rules) == 0 {
		err = reqerr.NewInvalidArgs("ExportFilter", "rules in filter should be empty")
		return
	}
	return
}

type ExportTsdbSpec struct {
	DestRepoName string            `json:"destRepoName"`
	SeriesName   string            `json:"series"`
	Tags         map[string]string `json:"tags"`
	Fields       map[string]string `json:"fields"`
	Timestamp    string            `json:"timestamp,omitempty"`
	Filter       *ExportFilter     `json:"filter,omitempty"`
}

func (s *ExportTsdbSpec) Validate() (err error) {
	if s.DestRepoName == "" {
		err = reqerr.NewInvalidArgs("ExportSpec", "dest repo name should not be empty")
		return
	}
	if s.SeriesName == "" {
		err = reqerr.NewInvalidArgs("ExportSpec", "series name should not be empty")
		return
	}
	if s.Filter == nil {
		return
	}
	return s.Filter.Validate()
}

type ExportMongoSpec struct {
	Host      string                 `json:"host"`
	DbName    string                 `json:"dbName"`
	CollName  string                 `json:"collName"`
	Mode      string                 `json:"mode"`
	UpdateKey []string               `json:"updateKey,omitempty"`
	Doc       map[string]interface{} `json:"doc"`
	Version   string                 `json:"version,omitempty"`
	Filter    *ExportFilter          `json:"filter,omitempty"`
}

func (s *ExportMongoSpec) Validate() (err error) {
	if s.Host == "" {
		err = reqerr.NewInvalidArgs("ExportSpec", "host should not be empty")
		return
	}
	if s.DbName == "" {
		err = reqerr.NewInvalidArgs("ExportSpec", "dbname should not be empty")
		return
	}
	if s.CollName == "" {
		err = reqerr.NewInvalidArgs("ExportSpec", "collection name should not be empty")
		return
	}
	if s.Mode != "UPSERT" && s.Mode != "INSERT" && s.Mode != "UPDATE" {
		err = reqerr.NewInvalidArgs("ExportSpec", fmt.Sprintf("invalid mode: %s, mode should be one of \"UPSERT\", \"INSERT\" and \"UPDATE\"", s.Mode))
		return
	}
	if s.Filter == nil {
		return
	}
	return s.Filter.Validate()
}

type ExportLogDBSpec struct {
	DestRepoName string                 `json:"destRepoName"`
	Doc          map[string]interface{} `json:"doc"`
	Filter       *ExportFilter          `json:"filter,omitempty"`
}

func (s *ExportLogDBSpec) Validate() (err error) {
	if s.DestRepoName == "" {
		err = reqerr.NewInvalidArgs("ExportSpec", "dest repo name should not be empty")
		return
	}
	if s.Filter == nil {
		return
	}
	return s.Filter.Validate()
}

type ExportKodoSpec struct {
	Bucket         string            `json:"bucket"`
	KeyPrefix      string            `json:"keyPrefix"`
	Fields         map[string]string `json:"fields"`
	RotateInterval int               `json:"rotateInterval,omitempty"`
	Email          string            `json:"email"`
	AccessKey      string            `json:"accessKey"`
	Format         string            `json:"format"`
	Compress       bool              `json:"compress"`
	Retention      int               `json:"retention"`
	Filter         *ExportFilter     `json:"filter,omitempty"`
}

func (s *ExportKodoSpec) Validate() (err error) {
	if s.Bucket == "" {
		err = reqerr.NewInvalidArgs("ExportSpec", "bucket should not be empty")
		return
	}
	if s.Filter == nil {
		return
	}
	return s.Filter.Validate()
}

type ExportHttpSpec struct {
	Host string `json:"host"`
	Uri  string `json:"uri"`
}

func (s *ExportHttpSpec) Validate() (err error) {
	if s.Host == "" {
		err = reqerr.NewInvalidArgs("ExportSpec", "host should not be empty")
		return
	}
	if s.Uri == "" {
		err = reqerr.NewInvalidArgs("ExportSpec", "uri should not be empty")
		return
	}
	return
}

type CreateExportInput struct {
	PipelineToken
	RepoName   string      `json:"-"`
	ExportName string      `json:"-"`
	Type       string      `json:"type"`
	Spec       interface{} `json:"spec"`
	Whence     string      `json:"whence,omitempty"`
}

type UpdateExportInput struct {
	PipelineToken
	RepoName   string      `json:"-"`
	ExportName string      `json:"-"`
	Spec       interface{} `json:"spec"`
}

func (e *UpdateExportInput) Validate() (err error) {
	if err = validateRepoName(e.RepoName); err != nil {
		return
	}
	if err = validateExportName(e.ExportName); err != nil {
		return
	}
	if e.Spec == nil {
		err = reqerr.NewInvalidArgs("ExportSpec", "spec should not be nil")
		return
	}
	switch e.Spec.(type) {
	case *ExportTsdbSpec, ExportTsdbSpec, *ExportMongoSpec, ExportMongoSpec,
		*ExportLogDBSpec, ExportLogDBSpec, *ExportKodoSpec, ExportKodoSpec,
		*ExportHttpSpec, ExportHttpSpec:
	default:
		return reqerr.NewInvalidArgs("ExportSpec", "spec Type not support")
	}
	return
}

func (e *CreateExportInput) Validate() (err error) {
	if err = validateRepoName(e.RepoName); err != nil {
		return
	}
	if err = validateExportName(e.ExportName); err != nil {
		return
	}
	if e.Spec == nil {
		err = reqerr.NewInvalidArgs("ExportSpec", "spec should not be nil")
		return
	}
	if e.Whence != "" && e.Whence != "oldest" && e.Whence != "newest" {
		err = reqerr.NewInvalidArgs("ExportSpec", "whence must be empty, \"oldest\" or \"newest\"")
		return
	}

	switch e.Spec.(type) {
	case *ExportTsdbSpec, ExportTsdbSpec:
		e.Type = "tsdb"
	case *ExportMongoSpec, ExportMongoSpec:
		e.Type = "mongo"
	case *ExportLogDBSpec, ExportLogDBSpec:
		e.Type = "logdb"
	case *ExportKodoSpec, ExportKodoSpec:
		e.Type = "kodo"
	case *ExportHttpSpec, ExportHttpSpec:
		e.Type = "http"
	default:
		return
	}

	vv, ok := e.Spec.(base.Validator)
	if !ok {
		err = reqerr.NewInvalidArgs("ExportSpec", "export spec cannot cast to validator")
		return
	}
	return vv.Validate()
}

type ExportDesc struct {
	Name   string                 `json:"name,omitempty"`
	Type   string                 `json:"type"`
	Spec   map[string]interface{} `json:"spec"`
	Whence string                 `json:"whence,omitempty"`
}

type GetExportInput struct {
	PipelineToken
	RepoName   string
	ExportName string
}

type GetExportOutput struct {
	ExportDesc
}

type ListExportsInput struct {
	PipelineToken
	RepoName string
}

type ListExportsOutput struct {
	Exports []ExportDesc `json:"exports"`
}

type DeleteExportInput struct {
	PipelineToken
	RepoName   string
	ExportName string
}

type VerifyTransformInput struct {
	PipelineToken
	Schema []RepoSchemaEntry `json:"schema"`
	Spec   *TransformSpec    `json:"spec"`
}

func (v *VerifyTransformInput) Validate() (err error) {
	if v.Schema == nil || len(v.Schema) == 0 {
		err = reqerr.NewInvalidArgs("Schema", "schema should not be empty")
		return
	}
	for _, item := range v.Schema {
		if err = item.Validate(); err != nil {
			return
		}
	}

	return v.Spec.Validate()
}

type VerifyTransformOutput struct {
	Schema []RepoSchemaEntry `json:"schema"`
}

type VerifyExportInput struct {
	PipelineToken
	Schema []RepoSchemaEntry `json:"schema"`
	Type   string            `json:"type"`
	Spec   interface{}       `json:"spec"`
	Whence string            `json:"whence,omitempty"`
}

func (v *VerifyExportInput) Validate() (err error) {
	if v.Schema == nil || len(v.Schema) == 0 {
		err = reqerr.NewInvalidArgs("VerifyExportSpec", "schema should not be empty")
		return
	}
	for _, item := range v.Schema {
		if err = item.Validate(); err != nil {
			return
		}
	}

	if v.Spec == nil {
		err = reqerr.NewInvalidArgs("ExportSpec", "spec should not be nil")
		return
	}

	if v.Whence != "" && v.Whence != "oldest" && v.Whence != "newest" {
		err = reqerr.NewInvalidArgs("ExportSpec", "whence must be empty, \"oldest\" or \"newest\"")
		return
	}

	switch v.Spec.(type) {
	case *ExportTsdbSpec, ExportTsdbSpec:
		v.Type = "tsdb"
	case *ExportMongoSpec, ExportMongoSpec:
		v.Type = "mongo"
	case *ExportLogDBSpec, ExportLogDBSpec:
		v.Type = "logdb"
	case *ExportKodoSpec, ExportKodoSpec:
		v.Type = "kodo"
	case *ExportHttpSpec, ExportHttpSpec:
		v.Type = "http"
	default:
		return
	}

	vv, ok := v.Spec.(base.Validator)
	if !ok {
		err = reqerr.NewInvalidArgs("ExportSpec", "export spec cannot cast to validator")
		return
	}
	return vv.Validate()
}

type KodoSourceSpec struct {
	Bucket      string   `json:"bucket"`
	KeyPrefixes []string `json:"keyPrefixes"`
	FileType    string   `json:"fileType"`
}

func (k *KodoSourceSpec) Validate() (err error) {
	if k.Bucket == "" {
		return reqerr.NewInvalidArgs("Bucket", fmt.Sprintf("bucket name should not be empty"))
	}
	if k.FileType == "" {
		return reqerr.NewInvalidArgs("FileType", fmt.Sprintf("fileType should not be empty"))
	}

	return
}

type HdfsSourceSpec struct {
	Paths    []string `json:"paths"`
	FileType string   `json:"fileType"`
}

func (h *HdfsSourceSpec) Validate() (err error) {
	if len(h.Paths) == 0 {
		return reqerr.NewInvalidArgs("Paths", fmt.Sprintf("paths should not be empty"))
	}
	for _, path := range h.Paths {
		if path == "" {
			return reqerr.NewInvalidArgs("Path", fmt.Sprintf("path in paths should not be empty"))
		}
	}
	if h.FileType == "" {
		return reqerr.NewInvalidArgs("FileType", fmt.Sprintf("fileType should not be empty"))
	}

	return
}

type RetrieveSchemaInput struct {
	PipelineToken
	Type string      `json:"type"`
	Spec interface{} `json:"spec"`
}

func (r *RetrieveSchemaInput) Validate() (err error) {
	switch r.Spec.(type) {
	case *KodoSourceSpec, KodoSourceSpec:
		r.Type = "kodo"
	case *HdfsSourceSpec, HdfsSourceSpec:
		r.Type = "hdfs"
	default:
		return
	}

	vv, ok := r.Spec.(base.Validator)
	if !ok {
		err = reqerr.NewInvalidArgs("Spec", "data source spec cannot cast to validator")
		return
	}
	return vv.Validate()
}

type RetrieveSchemaOutput struct {
	Schema []RepoSchemaEntry `json:"schema"`
}

type CreateDatasourceInput struct {
	PipelineToken
	DatasourceName string            `json:"-"`
	Region         string            `json:"region"`
	Type           string            `json:"type"`
	Spec           interface{}       `json:"spec"`
	Schema         []RepoSchemaEntry `json:"schema"`
}

func (c *CreateDatasourceInput) Validate() (err error) {
	if c.DatasourceName == "" {
		return reqerr.NewInvalidArgs("DatasourceName", fmt.Sprintf("datasource name should not be empty"))
	}
	if c.Type == "" {
		return reqerr.NewInvalidArgs("Type", fmt.Sprintf("type of datasource should not be empty"))
	}
	if len(c.Schema) == 0 {
		return reqerr.NewInvalidArgs("Schema", fmt.Sprintf("schema of datasource should not be empty"))
	}
	for _, schema := range c.Schema {
		if err = schema.Validate(); err != nil {
			return
		}
	}

	switch c.Spec.(type) {
	case *KodoSourceSpec, KodoSourceSpec:
		c.Type = "kodo"
	case *HdfsSourceSpec, HdfsSourceSpec:
		c.Type = "hdfs"
	default:
		return
	}

	vv, ok := c.Spec.(base.Validator)
	if !ok {
		err = reqerr.NewInvalidArgs("Spec", "data source spec cannot cast to validator")
		return
	}
	return vv.Validate()
}

type GetDatasourceInput struct {
	PipelineToken
	DatasourceName string
}

type GetDatasourceOutput struct {
	Region string            `json:"region"`
	Type   string            `json:"type"`
	Spec   interface{}       `json:"spec"`
	Schema []RepoSchemaEntry `json:"schema"`
}

type DatasourceDesc struct {
	Name   string            `json:"name"`
	Region string            `json:"region"`
	Type   string            `json:"type"`
	Spec   interface{}       `json:"spec"`
	Schema []RepoSchemaEntry `json:"schema"`
}

type ListDatasourcesOutput struct {
	Datasources []DatasourceDesc `json:"datasources"`
}

type DeleteDatasourceInput struct {
	PipelineToken
	DatasourceName string
}

type JobSrc struct {
	SrcName    string `json:"name"`
	FileFilter string `json:"fileFilter"`
	Type       string `json:"type"`
	TableName  string `json:"tableName"`
}

func (s *JobSrc) Validate() (err error) {
	if s.SrcName == "" {
		return reqerr.NewInvalidArgs("SrcName", fmt.Sprintf("source name should not be empty"))
	}
	if s.Type == "" {
		return reqerr.NewInvalidArgs("Type", fmt.Sprintf("source type should not be empty"))
	}
	if s.TableName == "" {
		return reqerr.NewInvalidArgs("TableName", fmt.Sprintf("table name should not be empty"))
	}

	return
}

type Computation struct {
	Code string `json:"code"`
	Type string `json:"type"`
}

func (c *Computation) Validate() (err error) {
	if c.Code == "" {
		return reqerr.NewInvalidArgs("Code", fmt.Sprintf("code in computation should not be empty"))
	}
	if c.Type == "" {
		return reqerr.NewInvalidArgs("Type", fmt.Sprintf("type in computation should not be empty"))
	}

	return
}

type JobSchedulerSpec struct {
	Crontab string `json:"crontab,omitempty"`
	Loop    string `json:"loop,omitempty"`
}

type JobScheduler struct {
	Type string            `json:"type"`
	Spec *JobSchedulerSpec `json:"spec,omitempty"`
}

type Param struct {
	Name    string `json:"name"`
	Default string `json:"default"`
}

type CreateJobInput struct {
	PipelineToken
	JobName     string        `json:"-"`
	Srcs        []JobSrc      `json:"srcs"`
	Computation Computation   `json:"computation"`
	Container   *Container    `json:"container,omitempty"`
	Scheduler   *JobScheduler `json:"scheduler,omitempty"`
	Params      []Param       `json:"params,omitempty"`
}

func (c *CreateJobInput) Validate() (err error) {
	if c.JobName == "" {
		return reqerr.NewInvalidArgs("JobName", fmt.Sprintf("job name should not be empty"))
	}
	if len(c.Srcs) == 0 {
		return reqerr.NewInvalidArgs("Srcs", fmt.Sprintf("must have at least one src inside the job srcs"))
	}
	for _, src := range c.Srcs {
		if err = src.Validate(); err != nil {
			return
		}
	}
	if err = c.Computation.Validate(); err != nil {
		return
	}

	return
}

type GetJobInput struct {
	PipelineToken
	JobName string
}

type GetJobOutput struct {
	Srcs        []JobSrc      `json:"srcs"`
	Computation Computation   `json:"computation"`
	Container   *Container    `json:"container,omitempty"`
	Scheduler   *JobScheduler `json:"scheduler,omitempty"`
	Params      []Param       `json:"params,omitempty"`
}

type JobDesc struct {
	Name        string        `json:"name"`
	Srcs        []JobSrc      `json:"srcs"`
	Computation Computation   `json:"computation"`
	Container   *Container    `json:"container,omitempty"`
	Scheduler   *JobScheduler `json:"scheduler,omitempty"`
	Params      []Param       `json:"params,omitempty"`
}

type ListJobsInput struct {
	PipelineToken
	SrcJobName        string
	SrcDatasourceName string
}

type ListJobsOutput struct {
	Jobs []JobDesc `json:"jobs"`
}

type DeleteJobInput struct {
	PipelineToken
	JobName string
}

type StartJobInput struct {
	PipelineToken
	JobName string  `json:"-"`
	Params  []Param `json:"params,omitempty"`
}

func (s *StartJobInput) Validate() (err error) {
	if s.JobName == "" {
		return reqerr.NewInvalidArgs("JobName", fmt.Sprintf("job name should not be empty"))
	}

	return
}

type StopJobInput struct {
	PipelineToken
	JobName string
}

type GetJobHistoryInput struct {
	PipelineToken
	JobName string
}

type JobHistory struct {
	RunId     int64  `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type GetJobHistoryOutput struct {
	Total   int64        `json:"total"`
	History []JobHistory `json:""`
}

type JobExportKodoSpec struct {
	Bucket      string   `json:"bucket"`
	KeyPrefix   string   `json:"keyPrefix"`
	Format      string   `json:"format"`
	Compression string   `json:"compression,omitempty"`
	Retention   int      `json:"retention"`
	PartitionBy []string `json:"partitionBy"`
	FileCount   int      `json:"fileCount"`
	Overwrite   bool     `json:"overwrite"`
}

func (e *JobExportKodoSpec) Validate() (err error) {
	if e.Bucket == "" {
		return reqerr.NewInvalidArgs("Bucket", fmt.Sprintf("bucket name should not be empty"))
	}
	if e.Format == "" {
		return reqerr.NewInvalidArgs("Format", fmt.Sprintf("format should not be empty"))
	}
	if e.FileCount <= 0 {
		return reqerr.NewInvalidArgs("FileCount", fmt.Sprintf("fileCount should be larger than 0"))
	}

	return
}

type CreateJobExportInput struct {
	PipelineToken
	JobName    string      `json:"-"`
	ExportName string      `json:"-"`
	Type       string      `json:"type"`
	Spec       interface{} `json:"spec"`
}

func (e *CreateJobExportInput) Validate() (err error) {
	if err = validateJobName(e.JobName); err != nil {
		return
	}
	if err = validateJobexportName(e.ExportName); err != nil {
		return
	}

	switch e.Spec.(type) {
	case *JobExportKodoSpec, JobExportKodoSpec:
		e.Type = "kodo"
	default:
		return
	}

	vv, ok := e.Spec.(base.Validator)
	if !ok {
		err = reqerr.NewInvalidArgs("JobExportSpec", "job export spec cannot cast to validator")
		return
	}
	return vv.Validate()
}

type GetJobExportInput struct {
	PipelineToken
	JobName    string
	ExportName string
}

type GetJobExportOutput struct {
	Type string      `json:"type"`
	Spec interface{} `json:"spec"`
}

type JobExportDesc struct {
	ExportName string      `json:"name"`
	Type       string      `json:"type"`
	Spec       interface{} `json:"spec"`
}

type ListJobExportsInput struct {
	PipelineToken
	JobName string
}

type ListJobExportsOutput struct {
	Exports []JobExportDesc `json:"exports"`
}

type DeleteJobExportInput struct {
	PipelineToken
	JobName    string
	ExportName string
}
