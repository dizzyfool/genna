package model

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

var reg = regexp.MustCompile(`/v\d+$`)

type CustomTypeMapping map[string]customType

func (c CustomTypeMapping) Add(pgType, goType, goImport string) {
	c[pgType] = customType{
		pgType:   pgType,
		goType:   goType,
		goImport: goImport,
	}
}

func (c CustomTypeMapping) Imports() []string {
	index := map[string]struct{}{}

	var result []string
	for _, customType := range c {
		if _, ok := index[customType.goImport]; ok {
			continue
		}

		if customType.goImport == "" {
			continue
		}

		result = append(result, customType.goImport)
		index[customType.goImport] = struct{}{}
	}

	return result
}

func (c CustomTypeMapping) Has(pgType string) bool {
	_, ok := c[pgType]
	return ok
}

func (c CustomTypeMapping) GoType(pgType string) (string, bool) {
	if customType, ok := c[pgType]; ok && customType.goType != "" {
		return customType.goType, true
	}

	return "", false
}

func (c CustomTypeMapping) GoImport(pgType string) (string, bool) {
	if customType, ok := c[pgType]; ok && customType.goImport != "" {
		return customType.goImport, true
	}

	return "", false
}

func ParseCustomTypes(raw []string) (CustomTypeMapping, error) {
	ctm := CustomTypeMapping{}

	for _, customType := range raw {
		pgType, goType, goImport, err := parseCustomType(customType)
		if err != nil {
			return nil, err
		}

		ctm.Add(pgType, goType, goImport)
	}

	return ctm, nil
}

func parseCustomType(raw string) (pgType, goType, goImport string, err error) {
	split := strings.SplitN(raw, ":", 2)
	if len(split) < 2 {
		err = fmt.Errorf("custom type mapping has invalid format (missing ':')")
		return
	}

	pgType = split[0]

	ind := strings.LastIndexByte(split[1], '.')
	if ind == -1 {
		goType = split[1]
		return
	}

	goImport, goType = split[1][:ind], split[1][ind:]

	if goType == "." || goImport == "" {
		err = fmt.Errorf("custom type mapping has invalid format (missing type or import)")
		return
	}

	if strings.Contains(goType, "/") {
		err = fmt.Errorf("type not found")
		return
	}

	ind = strings.LastIndexByte(goImport, '/')
	if ind == -1 {
		goType = fmt.Sprintf("%s%s", goImport, goType)
		return
	}

	base := path.Base(goImport)
	if reg.MatchString(goImport) {
		base = path.Base(path.Dir(goImport))
	}

	goType = fmt.Sprintf("%s%s", base, goType)

	return
}
