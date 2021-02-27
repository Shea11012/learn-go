package schema

import (
	"fmt"
	"strings"
)

func parseTag(tag string) map[string]string {
	tags := strings.Split(tag,",")
	var tagsMap = make(map[string]string,len(tags))
	for _,v := range tags {
		kv := strings.SplitN(v,":",2)
		var key,value string
		if len(kv) == 1 {
			key,value = kv[0],""
		}

		if len(kv) == 2 {
			key,value = kv[0],kv[1]
		}

		tagsMap[key] = value
	}

	return tagsMap
}

func Column(field *Field) string {
	var s strings.Builder
	s.WriteString(field.Name+" ")
	if v,ok := field.Tag["size"];ok {
		s.WriteString(fmt.Sprintf("%s(%s) ",field.Type,v))
		delete(field.Tag,"size")
	} else {
		s.WriteString(fmt.Sprintf("%s ",field.Type))
	}
	for k := range field.Tag {
		switch k {
		case "PRIMARY KEY":
			s.WriteString("PRIMARY KEY ")
		}
	}

	return s.String()
}