package mcron

import (
	"encoding/json"
	"log"
)

//统一传输协议 只响应 json格式 必须包含Action
//b := []byte(`{"Action":"register","Param":["Gomez","Morticia"]}`)
func jsonDecode(b []byte) map[string]interface{} {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		log.Println("非json数据：", err)
	}
	m := f.(map[string]interface{})
	return m
	//	log.Println(m["Action"], m["Param"])
	//	for k, v := range m {
	//		switch vv := v.(type) {
	//		case string:
	//			log.Println(k, "is string", vv)
	//		case int:
	//			log.Println(k, "is int", vv)
	//		case []interface{}:
	//			log.Println(k, "is an array:")
	//			for i, u := range vv {
	//				log.Println(i, u)
	//			}
	//		default:
	//			log.Println(k, "is of a type I don't know how to handle")
	//		}
	//	}
}

func jsonStr(key string, m map[string]interface{}) string {
	if nil == m[key] {
		return ""
	} else {
		return m[key].(string)
	}
}
func jsonArr(key string, m map[string]interface{}) []interface{} {
	v := m[key]
	switch vv := v.(type) {
	case []interface{}:
		return vv
	default:
		log.Println("is of a type I don't know how to handle")
	}
	return nil

}
