package gen

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"xbit/conf"

	"github.com/xbitgo/core/tools/tool_file"
	"github.com/xbitgo/core/tools/tool_str"

	"xbit/domian/gen/tpls"
	"xbit/domian/parser"
)

func (m *Manager) Pb(xstList map[string]parser.XST, otherList map[string]parser.XST) error {
	pkg := tool_str.LastPwdStr(m.AppPath)
	pkgPath := "proto/" + strings.Replace(m.AppPath, m.Project.RootPath()+"/", "", 1)
	pbGen := tpls.Pb{
		ProjectName: conf.Global.ProjectName,
		Import:      fmt.Sprintf("%s/apps/%s/domain/entity", conf.Global.ProjectName, m.AppName),
		Package:     pkg,
		PackagePath: pkgPath,
		MessageList: make([]tpls.Message, 0),
	}
	messageMap := make(map[string]tpls.Message)
	uniqMessage := make(map[string]struct{})
	for _, xst := range xstList {
		msgList := m.parseToPbMessage(xstList, otherList, xst, uniqMessage)
		for s, message := range msgList {
			messageMap[s] = message
		}
	}
	// 排序
	messageList := make([]tpls.Message, 0)
	for _, message := range messageMap {
		messageList = append(messageList, message)
	}
	sort.SliceStable(messageList, func(i, j int) bool {
		return messageList[i].Name < messageList[j].Name
	})
	pbGen.MessageList = messageList

	filename := fmt.Sprintf("%s/%s_%s_gen.proto", m.Tmpl.PbDir, pkg, "message")
	buf, err := pbGen.Execute()
	if err != nil {
		return err
	}
	buf = m.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		return err
	}

	m.pbConv(pbGen)
	return nil
}

func (m *Manager) pbConv(pbGen tpls.Pb) {
	// todo 扫描判断 如果存在 自定义转换方法 使用自定义转换方法处理 子字段转换
	buf, err := pbGen.ExecuteConv()
	if err != nil {
		fmt.Println(err)
		return
	}
	filename := fmt.Sprintf("%s/%s_converter_gen.go", m.Tmpl.ConvPbDir, "message")
	buf = m.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		fmt.Println(err)
		return
	}
}

var reMap, _ = regexp.Compile(`map\[(\w+)](\w+)`)

func (m *Manager) parseToPbMessage(xstList map[string]parser.XST, otherList map[string]parser.XST, xst parser.XST, uniqMessage map[string]struct{}) map[string]tpls.Message {
	uniqMessage[xst.Name] = struct{}{}
	messageList := make(map[string]tpls.Message, 0)

	msg := tpls.Message{
		Name:          xst.Name,
		MessageFields: make([]tpls.MessageField, 0),
	}
	// 排序
	fieldList := make([]parser.XField, 0)
	for _, field := range xst.FieldList {
		fieldList = append(fieldList, field)
	}
	sort.SliceStable(fieldList, func(i, j int) bool {
		return fieldList[i].Idx < fieldList[j].Idx
	})
	for i, field := range fieldList {
		name := tool_str.ToSnakeCase(field.Name)
		tpb := field.GetTag("pb")
		validate := field.GetTag("validate")
		if tpb == nil {
			continue
		}
		name = tpb.Name
		fe := tpls.MessageField{
			Type:     field.Type,
			Type2:    field.Type,
			Name:     name,
			Name2:    field.Name,
			NameSn:   tool_str.SnakeCaseToCamelCaseUF(name),
			Sort:     i + 1,
			JSONName: name,
			Label:    field.Comment,
		}
		if validate != nil {
			fe.Validate = fmt.Sprintf(`validate:"%s"`, validate.Name)
		}
		switch field.SType {
		case parser.STypeBasic:
			fe.Type = strings.TrimLeft(field.Type, "*")
			fe.Type, fe.ConvType = m.toTranPbBasicType(fe.Type)
		case parser.STypeStruct:
			fe.Type = strings.TrimLeft(field.Type, "*")
			if fe.Type == "time.Time" {
				if fe.Type == field.Type {
					fe.NoPoint = true
				}
				fe.Type = "string"
				fe.Type2 = "time.Time"
			} else {
				if fe.Type == field.Type {
					fe.NoPoint = true
				}
				if _, ok := xstList[fe.Type]; !ok {
					if cxs, ok := otherList[fe.Type]; ok {
						cList := m.parseToPbMessage(xstList, otherList, cxs, uniqMessage)
						for s, message := range cList {
							messageList[s] = message
						}
					}
				}
				fe.IsEntity = true
			}

		case parser.STypeSlice:
			fe.Type = strings.Replace(field.Type, "*", "", 1)
			if fe.Type == field.Type {
				fe.NoPoint = true
			}
			ident := strings.Replace(fe.Type, "[]", "", 1)

			fe.Type = strings.Replace(fe.Type, "[]byte", "bytes", -1)

			fe.Type = strings.Replace(fe.Type, "[]", "repeated ", 1)
			fe.Type, fe.ConvType = m.toTranPbBasicType(fe.Type)
			if ident == "time.Time" {
				fe.Type = strings.Replace(fe.Type, "time.Time", "string", 1)
			}
			fe.Type3 = ident
			if ident != "time.Time" && tool_str.UFirst(ident) {
				if _, ok := xstList[ident]; !ok {
					if cxs, ok := otherList[ident]; ok {
						if _, ok := uniqMessage[ident]; !ok {
							cList := m.parseToPbMessage(xstList, otherList, cxs, uniqMessage)
							for s, message := range cList {
								messageList[s] = message
							}
						}
					}
				}
				fe.IsMuEntity = true
				if fe.NoPoint {
					fe.EntityType = strings.Replace(field.Type, "[]", "[]entity.", 1)
					fe.PbType = strings.Replace(field.Type, "[]", "[]*pb.", 1)
				} else {
					fe.EntityType = strings.Replace(field.Type, "[]*", "[]*entity.", 1)
					fe.PbType = strings.Replace(field.Type, "[]*", "[]*pb.", 1)
				}
			}

		case parser.STypeMap:
			fe.Type = strings.Replace(field.Type, "*", "", 1)
			if fe.Type == field.Type {
				fe.NoPoint = true
			}
			rm := reMap.FindStringSubmatch(fe.Type)
			if len(rm) != 3 {
				continue
			}
			ident := rm[2]
			fe.Type = fmt.Sprintf("map<%s, %s>", rm[1], rm[2])
			fe.Type3 = rm[2]

			if ident == "time.Time" {
				fe.Type = strings.Replace(fe.Type, "time.Time", "string", 1)
			}
			fe.Type, _ = m.toTranPbBasicType(fe.Type)
			if ident != "time.Time" && tool_str.UFirst(ident) {
				if _, ok := xstList[ident]; !ok {
					if cxs, ok := otherList[ident]; ok {
						cList := m.parseToPbMessage(xstList, otherList, cxs, uniqMessage)
						for s, message := range cList {
							messageList[s] = message
						}
					}
				}
				fe.IsMuEntity = true
				if fe.NoPoint {
					fe.EntityType = fmt.Sprintf("map[%s]entity.%s", rm[1], rm[2])
					fe.PbType = fmt.Sprintf("map[%s]*pb.%s", rm[1], rm[2])
				} else {
					fe.EntityType = fmt.Sprintf("map[%s]*entity.%s", rm[1], rm[2])
					fe.PbType = fmt.Sprintf("map[%s]*pb.%s", rm[1], rm[2])
				}
			}
		}
		msg.MessageFields = append(msg.MessageFields, fe)
	}
	if len(msg.MessageFields) > 0 {
		messageList[msg.Name] = msg
	}
	return messageList
}

func (m *Manager) toTranPbBasicType(fType string) (nType string, convType bool) {
	switch fType {
	case "int":
		return "int64", true
	case "int8", "int16":
		return "int32", true
	case "float64":
		return "double", false
	case "float32":
		return "float", false
	}
	for k, v := range map[string]string{
		"<float64":         "<double",
		"float64>":         "double>",
		"float32>":         "float>",
		"<float32":         "<float",
		"repeated float64": "repeated double",
		"repeated float32": "repeated float",
	} {
		if strings.Contains(fType, k) {
			fType = strings.Replace(fType, k, v, 1)
		}
	}
	return fType, false
}
