package main

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm/schema"
)

type FieldInfo struct {
	FieldName string
	FieldType string
	Ignored   bool
	External  bool
}

func (fi FieldInfo) String() string {
	s := fmt.Sprintf("%s\t%s\t", fi.FieldName, fi.FieldType)
	if fi.Ignored {
		s += "ignored\t"
	}
	if fi.External {
		s += "external\t"
	}
	return s
}

type StructInfo struct {
	StructName    string
	FieldInfo     []*FieldInfo
	UniqueIndices map[string][]string
	PrimaryKeys   []string
	ColumnMap     map[string]string
}

func (si StructInfo) String() string {
	s := fmt.Sprintln(si.StructName)
	for _, fi := range si.FieldInfo {
		s += fmt.Sprintln(fi)
	}
	if len(si.UniqueIndices) > 0 {
		s += fmt.Sprintln("UniqueIndices:")
		for u, is := range si.UniqueIndices {
			s += fmt.Sprintln(u, strings.Join(is, "+"))
		}
	}
	if len(si.PrimaryKeys) > 0 {
		s += fmt.Sprintln("PrimaryKeys:", strings.Join(si.PrimaryKeys, "+"))
	}
	if len(si.ColumnMap) > 0 {
		s += fmt.Sprintln("Column names:")
		for fn, cn := range si.ColumnMap {
			s += fmt.Sprintln(fn, cn)
		}
	}
	return s
}

var testStructs = []string{
	"type Mouse struct {\n\tID           string         `gorm:\"primaryKey;comment:小鼠编号;\" json:\"id\"`\n\tYear         int            `json:\"year,omitempty\" gorm:\"index;comment:小鼠年份\"`\n\tNumber       int            `json:\"number,omitempty\" gorm:\"index;comment:小鼠编号数字部分\"`\n\tCreatedAt    time.Time      `json:\"created_at,omitempty\"`\n\tUpdatedAt    time.Time      `json:\"updated_at,omitempty\"`\n\tDeletedAt    gorm.DeletedAt `gorm:\"index\" json:\"-\"`\n\tGender       Gender         `json:\"gender\" swaggertype:\"primitive,string\" enums:\"male,female\" gorm:\"comment:性别\"`\n\tExperimentID string         `gorm:\"index;default:NULL;comment:实验项目号\" json:\"experiment_id,omitempty\"`\n\t// BirthdayOrArrivalDate is not necessarily the date of birth,\n\t// it can also be the arrival date depending on if the `WeekAgeOnArrival` is greater than 0\n\tBirthdayOrArrivalDate time.Time `json:\"boad\" gorm:\"column:boad;comment:出生日期/到货日期\"`\n\t// WeekAgeOnArrival is the week age of the mouse when it arrives at the facility,\n\t// if it's greater than 0 then the `BirthdayOrArrivalDate` means the arrival date\n\tWeekAgeOnArrival int                  `json:\"waoa,omitempty\" gorm:\"column:waoa;comment:到货周龄\"`\n\tProjectID        string               `gorm:\"index;comment:项目号\" json:\"project_id,omitempty\"`\n\tStrainID         uint                 `json:\"strain_id\" gorm:\"index;comment:品系号\"`\n\tStrain           Strain               `json:\"strain,omitempty\" gorm:\"foreignKey:StrainID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;\"`\n\tGenotype         *IdentifiedGenotypes `json:\"genotype,omitempty\" gorm:\"foreignKey:MouseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;\"`\n\tGeneration       int                  `json:\"generation\" gorm:\"comment:代数\"`\n\tFather           string               `json:\"father,omitempty\" gorm:\"comment:父亲\"`\n\tMother1          string               `json:\"mother1,omitempty\" gorm:\"comment:母亲1\"`\n\tMother2          string               `json:\"mother2,omitempty\" gorm:\"comment:母亲2\"`\n\tSupplier         string               `gorm:\"index;default:'Sironax';comment:供应商\" json:\"supplier,omitempty\"`\n\tStatus           Status               `gorm:\"index;comment:小鼠状态\" json:\"status\" swaggertype:\"primitive,string\" enums:\"feeding,genotyping,genotypechecking,unidentified,breedingwaiting,experimentwaiting,transferwaiting,newborn,breeding,experimentongoing,sacrificewaiting,dead,breedingfinished,endpoint,sacrificed,performingexternalexperiment,tissuecollection,genotypesunknown\"`\n\tCreatorID        uint                 `json:\"creator_id,omitempty\" gorm:\"comment:创建人\"`\n\tExternal         string               `json:\"external,omitempty\" gorm:\"index;default:NULL;comment:外部小鼠记录\"`\n\tRemarks          string               `json:\"remarks,omitempty\"`\n}",
	"type User struct {\n\tgorm.Model\n\tName     string   `json:\"name,omitempty\" gorm:\"uniqueIndex;comment:用户名\"`\n\tRealName string   `json:\"real_name,omitempty\" gorm:\"comment:用户真实姓名\"`\n\tPassword []byte   `json:\"-\" gorm:\"comment:用户密码哈希值\"`\n\tEmail    string   `json:\"email,omitempty\" gorm:\"uniqueIndex;comment:邮箱\"`\n\tPhone    string   `json:\"phone,omitempty\" gorm:\"uniqueIndex;comment:电话号码\"`\n\tRole     Role     `json:\"role,omitempty\" gorm:\"comment:用户角色\" swaggertype:\"primitive,string\" enums:\"watcher,identifier,experimenter,feeder,admin\"`\n\tToken    string   `json:\"token,omitempty\" gorm:\"-\"`\n\tConfig   []Config `json:\"config,omitempty\" gorm:\"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;\"`\n}\n",
}

var ns = schema.NamingStrategy{}

func squeeze(s string, c byte) string {
	var res []byte
	var flag bool
	for _, b := range []byte(s) {
		if b != c {
			res = append(res, b)
			if flag {
				flag = !flag
			}
			continue
		}
		if !flag {
			res = append(res, b)
			flag = !flag
		}
	}
	return string(res)
}

const testGoFile = `package main

// StrucA is a struct
type StructA struct {
	A int
	B int /*comments*/
	C string
}

/*
comments
comments
//
*/
type StructB struct {
	A uint // comments
	B int8 // comments
	C []byte
}
`

// RemoveComments removes non-code
func RemoveComments(content string) string {
	var res []string
	// remove comments enclosed by /**/
	if strings.Contains(content, "/*") && strings.Contains(content, "*/") {
		open := strings.Index(content, "/*")
		for open >= 0 {
			content = content[:open] + content[strings.Index(content, "*/")+2:]
			open = strings.Index(content, "/*")
		}
	}
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "//") {
			trimmed = strings.TrimSpace(trimmed[:strings.LastIndex(trimmed, "//")])
		}
		if trimmed == "" {
			continue
		}
		res = append(res, trimmed)
	}
	return strings.Join(res, "\n")
}

func FindStructBlock(structName string, fileContent string) string {
	pattern := fmt.Sprintf("type\\s+%s\\s+struct\\s+{", structName)
	if matched, _ := regexp.MatchString(pattern, fileContent); !matched {
		return ""
	}
	s := regexp.MustCompile(pattern).FindString(fileContent)
	height := 1
	start := strings.Index(fileContent, s)
	end := start + len(s)
	for _, b := range []byte(fileContent)[start+len(s):] {
		if b == '{' {
			height++
		} else if b == '}' {
			height--
		}
		if height == 0 {
			break
		}
		end++
	}
	return fileContent[start : end+1]
}

func ParseStructBlock(strStruct string) *StructInfo {
	fieldNameProcess := func(fn string) string {
		if strings.Contains(fn, ".") {
			fn = fn[strings.LastIndex(fn, ".")+1:]
		}
		return fn
	}
	lines := strings.Split(strStruct, "\n")
	structInfo := &StructInfo{
		FieldInfo:     []*FieldInfo{},
		UniqueIndices: make(map[string][]string),
		PrimaryKeys:   []string{},
		ColumnMap:     make(map[string]string),
	}
	regexTag := regexp.MustCompile("`.*`")
	for _, line := range lines {
		if matched, _ := regexp.MatchString(`type\s+\w+\s+struct`, line); matched {
			structName := strings.TrimSpace(strings.TrimLeft(strings.TrimRight(regexp.MustCompile(`type\s+\w+\s+struct`).FindString(line), "struct"), "type"))
			structInfo.StructName = structName
			continue
		}
		var tag, trimmed string
		if tag = regexTag.FindString(line); tag != "" {
			trimmed = strings.TrimSpace(strings.ReplaceAll(line, tag, ""))
		} else {
			trimmed = strings.TrimSpace(line)
		}
		if trimmed == "" || strings.HasSuffix(trimmed, "}") {
			continue
		}
		splits := strings.Split(squeeze(trimmed, ' '), " ")
		var fieldname, fieldtype string
		if len(splits) == 2 {
			fieldname, fieldtype = splits[0], splits[1]
		} else if len(splits) == 1 {
			fieldname, fieldtype = fieldNameProcess(splits[0]), splits[0]
		}
		fieldInfo := &FieldInfo{
			FieldName: fieldname,
			FieldType: fieldtype,
		}
		if matched, _ := regexp.MatchString(`gorm:"-[^"]*"`, tag); matched {
			fieldInfo.Ignored = true
			continue
		}
		if matched, _ := regexp.MatchString(`gorm:"[^"]+"`, tag); matched {
			structInfo.ColumnMap[fieldname] = ns.ColumnName(structInfo.StructName, fieldname)
			gormInfo := strings.TrimSpace(regexp.MustCompile(`gorm:"[^"]+"`).FindString(tag))
			gormInfo = strings.TrimSpace(gormInfo[6 : len(gormInfo)-1])
			gormInfos := strings.Split(gormInfo, ";")
			for _, gi := range gormInfos {
				if trmd := strings.ToLower(strings.TrimSpace(gi)); trmd != "" {
					if strings.Contains(trmd, ":") {
						splts := strings.Split(trmd, ":")
						k, v := splts[0], splts[1]
						if k == "uniqueindex" {
							if _, ok := structInfo.UniqueIndices[v]; !ok {
								structInfo.UniqueIndices[v] = []string{}
							}
							structInfo.UniqueIndices[v] = append(structInfo.UniqueIndices[v], fieldname)
						} else if k == "column" {
							structInfo.ColumnMap[fieldname] = v
						} else if k == "foreignkey" {
							fieldInfo.External = true
							delete(structInfo.ColumnMap, fieldname)
							break
						} else if k == "embeddedprefix" {
							// TODO: complete this section
						}
					} else {
						if trmd == "uniqueindex" {
							structInfo.UniqueIndices[fieldname] = []string{fieldname}
						} else if trmd == "primarykey" {
							structInfo.PrimaryKeys = append(structInfo.PrimaryKeys, fieldname)
						}
					}
				}
			}
		}
		structInfo.FieldInfo = append(structInfo.FieldInfo, fieldInfo)
	}
	return structInfo
}

func main() {
	for _, testStruct := range testStructs {
		structInfo := ParseStructBlock(RemoveComments(testStruct))
		fmt.Println(*structInfo)
	}

	sa := FindStructBlock("StructA", RemoveComments(testGoFile))
	fmt.Println(sa)
	fmt.Println()
	structA := ParseStructBlock(sa)
	fmt.Println(*structA)
	sb := FindStructBlock("StructB", RemoveComments(testGoFile))
	fmt.Println(sb)
	fmt.Println()
	structB := ParseStructBlock(sb)
	fmt.Println(*structB)
}
