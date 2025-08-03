package generator

import (
	"fmt"
	"strings"

	"github.com/arturh85/pocketbase-go-generator/internal/cmd"
	"github.com/arturh85/pocketbase-go-generator/internal/pocketbase_api"
	"github.com/iancoleman/strcase"
)

type InterfacePropertyType int

const (
	IptString = iota
	IptNumber
	IptBoolean
	IptJson
	IptFile
	IptEnum
	IptRelation
	IptDate
)

type InterfaceProperty struct {
	Name           string
	CollectionName string
	Optional       bool
	Type           InterfacePropertyType
	IsArray        bool
	Data           interface{}
}

type CollectionWithProperties struct {
	Collection *pocketbase_api.Collection
	Properties []*InterfaceProperty
}

type propertyFlags struct {
	relationAsString bool
	forceOptional    bool
}

func GetInterfacePropertyType(typeName string) InterfacePropertyType {
	switch typeName {
	case "number":
		return IptNumber
	case "bool":
		return IptBoolean
	case "select":
		return IptEnum
	case "json":
		return IptJson
	case "file":
		return IptFile
	case "relation":
		return IptRelation
	case "date":
		return IptDate
	case "autodate":
		return IptDate
	default:
		return IptString
	}
}

func (propertyType InterfacePropertyType) String() string {
	switch propertyType {
	case IptString:
		return "String"
	case IptNumber:
		return "Number"
	case IptBoolean:
		return "Boolean"
	case IptEnum:
		return "Enum"
	case IptJson:
		return "Json"
	case IptFile:
		return "File"
	case IptRelation:
		return "Relation"
	}

	return "Unknown"
}

func (property InterfaceProperty) String() string {
	var data = []string{
		property.Type.String(),
	}

	if property.Optional {
		data = append(data, "Optional")
	}

	if property.IsArray {
		data = append(data, "Array")
	}

	if property.Type == IptRelation {
		relationTo, ok := property.Data.(string)
		if !ok {
			relationTo = "unknown (object)"
		}

		data = append(data, fmt.Sprintf("Relation to %s", relationTo))
	}

	if property.Type == IptEnum {
		enumData := property.Data.([]string)

		data = append(data, fmt.Sprintf("Enum Data [%s]", strings.Join(enumData, ", ")))
	}

	return fmt.Sprintf("%s (%s)", property.Name, strings.Join(data, ", "))
}

func (property InterfaceProperty) GetGoProperty(generatorFlags *cmd.GeneratorFlags, flags propertyFlags) string {
	// IntValue int `json:"intValue"`
	return fmt.Sprintf("%s %s `json:\"%s\"`", strcase.ToCamel(property.getGoName(generatorFlags, flags)), property.getGoTypeWithArray(flags), property.getGoName(generatorFlags, flags))
}

/*
example getter:

	func (a *Article) Title() string {
	    return a.GetString("title")
	}
*/
func (property InterfaceProperty) GetGoRecordGetter(generatorFlags *cmd.GeneratorFlags, flags propertyFlags) string {
	intGetter := ""

	if property.Type == IptNumber {
		intGetter = fmt.Sprintf("func (a *%sRecord) %s() %s {\n    return a.%s(\"%s\")\n}\n\n",
			strcase.ToCamel(property.CollectionName),
			strcase.ToCamel(property.Name)+"Int",
			"int",
			"GetInt",
			property.getGoName(generatorFlags, flags),
		)
	}

	return fmt.Sprintf("%sfunc (a *%sRecord) %s() %s {\n    return a.%s(\"%s\")\n}\n",
		intGetter,
		strcase.ToCamel(property.CollectionName),
		strcase.ToCamel(property.Name),
		property.getGoRecordType(flags),
		property.getPocketbaseGetter(flags),
		property.getGoName(generatorFlags, flags),
	)
}
func (property InterfaceProperty) GetGoRecordExpandRelation(generatorFlags *cmd.GeneratorFlags, flags propertyFlags) string {

	// if errs := app.ExpandRecord(wjob.Record, []string{collections.WorkerJobsFields.Job}, nil); len(errs) > 0 {
	// 	logrus.Error(fmt.Sprintf("Failed to expand job for worker job %s: %v", wjob.Id(), errs))
	// 	sentry.CaptureException(err)
	// 	continue
	// }
	// job := collections.Jobs_Wrap(wjob.ExpandedOne(collections.WorkerJobsFields.Job))

	relationName := property.Data.(string)
	return fmt.Sprintf(`		
		func (a *%sRecord) Expand%s(app core.App) (%s, error) {
			if errs := app.ExpandRecord(a.Record, []string{"%s"}, nil); len(errs) > 0 {
				return nil, errs["%s"]
			}
			record := a.ExpandedOne("%s")
			if record == nil { return nil, nil }
			return %s_Wrap(record), nil
		}
		`,
		strcase.ToCamel(property.CollectionName),
		strcase.ToCamel(property.Name),
		"*"+strcase.ToCamel(relationName)+"Record",

		property.getGoName(generatorFlags, flags),
		property.getGoName(generatorFlags, flags),
		property.getGoName(generatorFlags, flags),
		strcase.ToCamel(relationName),
	)
}

var reserved = []string{"type", "var", "func", "map", "struct"}

/*
example setter:

	func (a *Article) SetTitle(title string) {
	    a.Set("title", title)
	}
*/

func (property InterfaceProperty) GetGoRecordSetter(generatorFlags *cmd.GeneratorFlags, flags propertyFlags) string {
	var validGoName = strcase.ToLowerCamel(property.getGoName(generatorFlags, flags))
	for _, word := range reserved {
		if validGoName == word {
			validGoName = "_" + validGoName
			break
		}
	}

	if property.Name == "created" || property.Name == "updated" {
		return ""
	}

	return fmt.Sprintf("func (a *%sRecord) Set%s(%s %s) {\n    a.Set(\"%s\", %s)\n}\n",
		strcase.ToCamel(property.CollectionName),
		strcase.ToCamel(property.Name),
		validGoName,
		property.getGoRecordType(flags),
		property.getGoName(generatorFlags, flags),
		validGoName,
	)
}

func (property InterfaceProperty) getGoType(flags propertyFlags) string {
	switch property.Type {
	case IptNumber:
		if property.Optional {
			return "*float32"
		} else {
			return "float32"
		}
	case IptBoolean:
		if property.Optional {
			return "*bool"
		} else {
			return "bool"
		}
	case IptJson:
		if property.Optional {
			return "*map[string]interface{}"
		} else {
			return "map[string]interface{}"
		}
	case IptEnum:
		return strcase.ToCamel(fmt.Sprintf("%s_%s_%s", property.CollectionName, property.Name, "options"))
	case IptRelation:
		if flags.relationAsString {
			return "string"
		}

		relationTo, ok := property.Data.(string)
		if !ok {
			return "map[string]interface{}"
		} else {
			if property.Optional {
				return "*" + strcase.ToCamel(relationTo) + "Struct"
			} else {
				return strcase.ToCamel(relationTo) + "Struct"
			}
		}
	default:
		if property.Optional {
			return "*string"
		} else {
			return "string"
		}
	}
}
func (property InterfaceProperty) getPocketbaseGetter(flags propertyFlags) string {
	if property.IsArray {
		return "GetStringSlice"
	}
	switch property.Type {
	case IptNumber:
		return "GetFloat"
	case IptBoolean:
		return "GetBool"
	case IptJson:
		return "Get"
	case IptEnum:
		return "GetString"
	case IptRelation:
		return "GetString"
	case IptDate:
		return "GetDateTime"
	default:
		return "GetString"
	}
}
func (property InterfaceProperty) getGoRecordType(flags propertyFlags) string {
	if property.IsArray {
		return "[]string"
	}
	switch property.Type {
	case IptNumber:
		return "float64"
	case IptBoolean:
		return "bool"
	case IptJson:
		return "any"
	case IptEnum:
		return "string"
	case IptRelation:
		return "string"
	case IptDate:
		return "types.DateTime"
	default:
		return "string"
	}
}

func (property InterfaceProperty) getGoTypeWithArray(flags propertyFlags) string {
	tsType := property.getGoType(flags)

	if property.IsArray {
		if strings.HasPrefix(tsType, "*") {
			tsType = tsType[1:]
		}
		if property.Optional {
			return fmt.Sprintf("*[]%s", tsType)
		} else {
			return fmt.Sprintf("[]%s", tsType)
		}
	}

	return tsType
}

func (property InterfaceProperty) getGoName(generatorFlags *cmd.GeneratorFlags, flags propertyFlags) string {
	if property.Optional && generatorFlags.MakeNonRequiredOptional || flags.forceOptional {
		return fmt.Sprintf("%s", property.Name)
	}

	return property.Name
}

func (collection CollectionWithProperties) GetGoCollectionEntry(generatorFlags *cmd.GeneratorFlags) string {
	return fmt.Sprintf("    Collection%s = \"%s\"", strcase.ToCamel(collection.Collection.Name), collection.Collection.Name)
}

func (collection CollectionWithProperties) GetGoCollectionHelperFuncs(generatorFlags *cmd.GeneratorFlags) string {
	template := `
func $$$_Wrap(record *core.Record) *$$$Record {
	typedRecord := &$$$Record{}
	typedRecord.SetProxyRecord(record)
	return typedRecord
}

func $$$_New(app core.App) (*$$$Record, error) {
	c, err := app.FindCollectionByNameOrId(Collection$$$)
	if err != nil { return nil, err }
	return $$$_Wrap(core.NewRecord(c)), nil
}

func $$$_FindRecordById(app core.App, recordId string, optFilters ...func(q *dbx.SelectQuery) error) (*$$$Record, error) {
	_record, err := app.FindRecordById(Collection$$$, recordId, optFilters...)
	if err != nil { return nil, err }
	return $$$_Wrap(_record), nil
}

func $$$_FindFirstRecordByData(app core.App, key string, value any) (*$$$Record, error) {
	_record, err := app.FindFirstRecordByData(Collection$$$, key, value)
	if err != nil { return nil, err }
	return $$$_Wrap(_record), nil
}

func $$$_FindRecordsByFilter(app core.App,
	filter string,
	sort string,
	limit int,
	offset int,
	params ...dbx.Params) ([]*$$$Record, error) {
	_records, err := app.FindRecordsByFilter(Collection$$$, filter, sort, limit, offset, params...)
	if err != nil { return nil, err }
	records := make([]*$$$Record, len(_records))
	for i, _record := range _records { records[i] = $$$_Wrap(_record) }
	return records, err
}
	`
	return strings.ReplaceAll(template, "$$$", strcase.ToCamel(collection.Collection.Name))
}

func (collection CollectionWithProperties) GetGoRecord(generatorFlags *cmd.GeneratorFlags) string {
	properties := make([]string, len(collection.Properties))
	var additionalTypes []string

	for i, property := range collection.Properties {
		properties[i] = fmt.Sprintf("%s\n\n%s",
			property.GetGoRecordGetter(generatorFlags, propertyFlags{forceOptional: false, relationAsString: true}),
			property.GetGoRecordSetter(generatorFlags, propertyFlags{forceOptional: false, relationAsString: true}))

		if property.Type == IptRelation {
			properties[i] += property.GetGoRecordExpandRelation(generatorFlags, propertyFlags{forceOptional: false, relationAsString: true})
		}
	}

	prefix := strings.Join(additionalTypes, "\n\n")

	if prefix != "" {
		prefix += "\n\n"
	}

	var publicExportStruct = `
	func (a *$$$Record) PublicExportStruct() $$$Struct {
		bytes, _ := json.Marshal(a.PublicExport())
		var record = $$$Struct{}
		_ = json.Unmarshal(bytes, &record)
		return record
	}
	`

	return fmt.Sprintf("%s\nvar _ core.RecordProxy = (*%sRecord)(nil)\n\ntype %sRecord struct {\n    core.BaseRecordProxy\n}\n\n%s\n\n%s\n\n",
		prefix,
		strcase.ToCamel(collection.Collection.Name),
		strcase.ToCamel(collection.Collection.Name),
		strings.Join(properties, "\n"),
		strings.ReplaceAll(publicExportStruct, "$$$", strcase.ToCamel(collection.Collection.Name)),
	)
}

func (collection CollectionWithProperties) GetGoStruct(generatorFlags *cmd.GeneratorFlags) string {
	properties := make([]string, len(collection.Properties))
	var additionalTypes []string
	var expandedRelations []string

	fieldNames := make([]string, len(collection.Properties))
	fieldNameValues := make([]string, len(collection.Properties))

	for i, property := range collection.Properties {
		fieldNames[i] = strcase.ToCamel(property.Name)
		fieldNameValues[i] = fmt.Sprintf("%s: \"%s\"", fieldNames[i], property.Name)
		properties[i] = fmt.Sprintf("    %s;", property.GetGoProperty(generatorFlags, propertyFlags{forceOptional: false, relationAsString: true}))

		if property.Type == IptEnum {
			additionalTypes = append(additionalTypes, property.getGoEnum())
		}

		if property.Type == IptRelation {
			expandedRelations = append(expandedRelations, fmt.Sprintf("    %s;", property.GetGoProperty(generatorFlags, propertyFlags{forceOptional: true, relationAsString: false})))
		}
	}

	if len(expandedRelations) > 0 {
		// expandedRelations = append(expandedRelations, "    [key: string]: unknown;")

		expandedType := fmt.Sprintf("type %sExpanded struct {\n%s\n}", strcase.ToCamel(collection.Collection.Name), strings.Join(expandedRelations, "\n"))

		additionalTypes = append(additionalTypes, expandedType)

		expandedLine := fmt.Sprintf("    Expand %sExpanded `json:\"expand\"`", strcase.ToCamel(collection.Collection.Name))

		properties = append([]string{expandedLine}, properties...)
	} else {
		// expandedLine := "    expand?: { [key: string]: unknown; };"

		// properties = append([]string{expandedLine}, properties...)
	}

	prefix := strings.Join(additionalTypes, "\n\n")

	if prefix != "" {
		prefix += "\n\n"
	}

	var fieldsInfo = fmt.Sprintf("var %sFields = struct {\n    %s string\n}{\n%s,\n}", strcase.ToCamel(collection.Collection.Name), strings.Join(fieldNames, ", "), strings.Join(fieldNameValues, ",\n"))
	return fmt.Sprintf("%stype %sStruct struct {\n%s\n}\n\n%s", prefix, strcase.ToCamel(collection.Collection.Name), strings.Join(properties, "\n"), fieldsInfo)
}

func (property InterfaceProperty) getGoEnum() string {
	if property.Type != IptEnum {
		return ""
	}

	enumData := property.Data.([]string)
	enumName := strcase.ToCamel(fmt.Sprintf("%s_%s_%s", property.CollectionName, property.Name, "options"))

	enumList := make([]string, len(enumData))

	for i, enum := range enumData {
		enumList[i] = fmt.Sprintf("    %s %s = \"%s\"", enumName+"_"+strcase.ToCamel(enum), enumName, enum)
	}

	return fmt.Sprintf("type %s string\nconst (\n%s\n)", enumName, strings.Join(enumList, "\n"))
}
