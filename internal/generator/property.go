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
				return "*" + strcase.ToCamel(relationTo)
			} else {
				return strcase.ToCamel(relationTo)
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

func (collection CollectionWithProperties) GetGoInterface(generatorFlags *cmd.GeneratorFlags) string {
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
	return fmt.Sprintf("%stype %s struct {\n%s\n}\n\n%s", prefix, strcase.ToCamel(collection.Collection.Name), strings.Join(properties, "\n"), fieldsInfo)
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
