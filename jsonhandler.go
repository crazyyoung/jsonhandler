//Copyright 2017 crazyyoung
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
package jsonhandler

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

//枚举json数据类型
//Enumeration of JSON data types
type jsonType int

const (
	JArray jsonType = iota
	JObject
	JBool
	JNumber
	JString
	JNull
)

//JObjectValue 内部存储json object 使用的数据类型
//JObjectValue Internal type for json object
type JObjectValue map[string]*JsonNode

//JArrayValue 内部存储json array 使用的数据类型
//JArrayValue Internal type for json array
type JArrayValue []*JsonNode

//JsonNode 内部存储一个json值使用的数据类型
//JsonNode Internal type for json elements
type JsonNode struct {
	nodeType jsonType    //内部存储的类型
	value    interface{} //实际保存数据的空间
}

//typeToJsonType 将go的数类型转换成内部描述json值的数据类型
//typeToJsonType from go type to internal type
func typeToJsonType(f interface{}) (jsonType, error) {
	switch f.(type) {
	case ([]interface{}): //数组
		return JArray, nil
	case (map[string]interface{}): //map
		return JObject, nil
	case bool:
		return JBool, nil
	case float64:
		return JNumber, nil
	case string:
		return JString, nil
	case nil:
		return JNull, nil
	default:
		return JNull, errors.New("wrong type")
	}
}

//supportParse 利用encoding/json包处理的结果，将json字符串生成为一个内部定义的数据结构
//supportParse use encoding/json to Unmarshal json string
func supportUnmarshal(f interface{}) (*JsonNode, error) {
	var returnNode JsonNode
	var err error
	switch f.(type) {
	case ([]interface{}): //数组
		returnNode.nodeType = JArray
		vList := make(JArrayValue, 0)
		for _, v := range f.([]interface{}) {
			var getNode *JsonNode
			getNode, err = supportUnmarshal(v)
			if err != nil {
				return nil, err
			}
			vList = append(vList, getNode)
		}
		returnNode.value = vList
	case (map[string]interface{}): //map
		returnNode.nodeType = JObject
		vMap := make(JObjectValue)
		for k, v := range f.(map[string]interface{}) {
			var getNode *JsonNode
			getNode, err = supportUnmarshal(v)
			if err != nil {
				return nil, err
			}
			vMap[k] = getNode
		}
		returnNode.value = vMap
	case bool:
		returnNode.nodeType = JBool
		returnNode.value = f.(bool)
	case float64:
		returnNode.nodeType = JNumber
		returnNode.value = f.(float64)
	case string:
		returnNode.nodeType = JString
		returnNode.value = f.(string)
	case nil:
		returnNode.nodeType = JNull
		returnNode.value = nil
	default:
		return nil, errors.New("wrong type " + reflect.TypeOf(f).String())
	}
	return &returnNode, err
}

//Copy 复制以一个json值信息
//Copy Copy a json element
func (node *JsonNode) Copy(copyFrom *JsonNode) (err error) {
	node.nodeType = copyFrom.nodeType
	node.value = copyFrom.value
	return nil
}

//Unmarshal 利用encoding/json包处理的结果，调用supportParse生成为一个内部定义的数据结构
//Unmarshal use encoding/json to Unmarshal json string
func (node *JsonNode) Unmarshal(body []byte) (err error) {
	var f interface{}
	var returnNode *JsonNode
	err = json.Unmarshal(body, &f)
	if err != nil {
		return err
	}
	returnNode, err = supportUnmarshal(f)
	if err != nil {
		return err
	}
	node.Copy(returnNode)
	return err
}

//NewJsonNode 新建存储json信息的结构
//NewJsonNode new a null json element
func NewJsonNode() (*JsonNode, error) {
	var returnNode JsonNode
	returnNode.nodeType = JNull
	returnNode.value = nil
	return &returnNode, nil
}

//Bool 按布尔型返回信息
//Bool return bool type json element
func (node *JsonNode) Bool() (bool, error) {
	if !node.IsBool() {
		return false, errors.New("node is not bool type")
	}
	return node.value.(bool), nil
}

//GetBool 按布尔型返回信息
//GetBool return bool type json element
func (node *JsonNode) GetBool(args ...interface{}) (bool, error) {
	r, e := node.Get(args[0:]...)
	if !r.IsBool() {
		return false, errors.New("node is not bool type")
	}
	return r.value.(bool), e
}

//Number 按数字返回信息
//Number return number type json element
func (node *JsonNode) Number() (float64, error) {
	if !node.IsNumber() {
		return 0, errors.New("node is not number type")
	}
	return node.value.(float64), nil
}

//GetNumber 按数字返回信息
//GetNumber return number type json element
func (node *JsonNode) GetNumber(args ...interface{}) (float64, error) {
	r, e := node.Get(args[0:]...)
	if !r.IsNumber() {
		return 0, errors.New("node is not number type")
	}
	return r.value.(float64), e
}

//String 按字符串返回信息
//String return string type json element
func (node *JsonNode) String() (string, error) {
	if !node.IsString() {
		return "", errors.New("node is not string type")
	}
	return node.value.(string), nil
}

//GetString 按字符串返回信息
//GetString return string type json element
func (node *JsonNode) GetString(args ...interface{}) (string, error) {
	r, e := node.Get(args[0:]...)
	if !r.IsString() {
		return "", errors.New("node is not string type")
	}
	return r.value.(string), e
}

//Null 按Null返回信息
//Null return null type json element
func (node *JsonNode) Null() (interface{}, error) {
	if !node.IsNull() {
		return nil, errors.New("node is not null type")
	}
	return nil, nil
}

//GetNull 按Null返回信息
//GetNull return null type json element
func (node *JsonNode) GetNull(args ...interface{}) (interface{}, error) {
	r, e := node.Get(args[0:]...)
	if !r.IsNull() {
		return nil, errors.New("node is not null type")
	}
	return nil, e
}

func (node *JsonNode) getNextNode(arg interface{}) (returnNode *JsonNode, err error) {
	switch arg.(type) {
	case string:
		if node.nodeType == JObject {
			return (node.value.(JObjectValue))[arg.(string)], err
		} else {
			return nil, errors.New("is not a json object")
		}
	case int:
		if node.nodeType == JArray {
			return (node.value.(JArrayValue))[arg.(int)], err
		} else {
			return nil, errors.New("is not a json array")
		}
	default:
		return nil, errors.New("key is not an acceptable type.type is " + reflect.TypeOf(arg).String())
	}
}

//Get 获取一个json信息
//Get get a json element
func (node *JsonNode) Get(args ...interface{}) (returnNode *JsonNode, err error) {
	argsLen := len(args)
	if argsLen <= 0 {
		returnNode = nil
		err = errors.New("wrong type")
	} else if argsLen == 1 {
		arg := args[0]
		returnNode, err = node.getNextNode(arg)
	} else {
		var firstNode *JsonNode
		arg := args[0]
		firstNode, err = node.getNextNode(arg)
		if err != nil {
			return nil, err
		}
		return firstNode.Get(args[1:]...)
	}
	return returnNode, err
}
func (node *JsonNode) getOrAddNextNode(arg interface{}) (returnNode *JsonNode, err error) {
	switch arg.(type) {
	case string:
		if node.nodeType == JObject {
			_, ok := (node.value.(JObjectValue))[arg.(string)]
			if !ok {
				var newNode *JsonNode
				newNode, err = NewJsonNode()
				if err != nil {
					return nil, err
				}
				(node.value.(JObjectValue))[arg.(string)] = newNode
			}
			return (node.value.(JObjectValue))[arg.(string)], err
		} else {
			node.nodeType = JObject
			node.value = make(JObjectValue)
			return node.getOrAddNextNode(arg)
		}
	case int:
		if node.nodeType == JArray {
			arrLen := len((node.value.(JArrayValue)))
			if arg.(int) < 0 {
				return nil, errors.New("index can not < 0")
			} else if arg.(int) >= arrLen {
				for i := arrLen; i <= arg.(int); i++ {
					var newNode *JsonNode
					newNode, err = NewJsonNode()
					if err != nil {
						return nil, err
					}
					node.value = append(node.value.(JArrayValue), newNode)
				}
			}
			return (node.value.(JArrayValue))[arg.(int)], err
		} else {
			node.nodeType = JArray
			node.value = make(JArrayValue, 0)
			return node.getOrAddNextNode(arg)
		}
	default:
		return nil, errors.New("key is not an acceptable type.type is " + reflect.TypeOf(arg).String())
	}
}

//Set 设置一个json信息
//Set set a json element
func (node *JsonNode) Set(value interface{}, args ...interface{}) (returnNode *JsonNode, err error) {
	//整数转换成浮点数用
	if reflect.TypeOf(value).Kind() == reflect.Int {
		value = float64(value.(int))
	}
	argsLen := len(args)
	if argsLen <= 0 {
		switch value.(type) {
		case *JsonNode:
			node.nodeType = value.(*JsonNode).nodeType
			node.value = value.(*JsonNode).value
		default:
			node.value = value
			node.nodeType, err = typeToJsonType(value)
		}
		returnNode = node
	} else if argsLen == 1 {
		var firstNode *JsonNode
		arg := args[0]
		firstNode, err = node.getOrAddNextNode(arg)
		return firstNode.Set(value)
	} else {
		var firstNode *JsonNode
		arg := args[0]
		firstNode, err = node.getOrAddNextNode(arg)
		if err != nil {
			return nil, err
		}
		return firstNode.Set(value, args[1:]...)
	}
	return returnNode, err
}

//Delete 删除一个json信息
//Delete delete a json element
func (node *JsonNode) Delete(args ...interface{}) (ok bool, err error) {
	argsLen := len(args)
	if argsLen <= 0 {
		node.nodeType = JNull
		node.value = nil
		node = nil
		return false, errors.New("can not delete without parameter")
	} else if argsLen == 1 {
		switch node.nodeType {
		case JObject:
			if reflect.TypeOf(args[0]).Kind() == reflect.String {
				var key string
				key = args[0].(string)
				if _, ok := node.value.(JObjectValue)[key]; ok {
					delete(node.value.(JObjectValue), key)
					return true, nil
				}
			}
		case JArray:
			if reflect.TypeOf(args[0]).Kind() == reflect.Int {
				var index int
				index = args[0].(int)
				if index >= 0 && index < len(node.value.(JArrayValue)) {
					node.value = append(node.value.(JArrayValue)[:index], node.value.(JArrayValue)[index+1:]...)
				}
			}
		}
		return false, errors.New("there is no json info here")
	} else {
		r, err := node.Get(args[0 : argsLen-1]...)
		if err != nil {
			return false, errors.New("there is no json info here")
		}
		return r.Delete(args[argsLen-1])
	}
}

//IsArray 类型检测
//IsArray test if is a array json element
func (node *JsonNode) IsArray() bool {
	if node.nodeType == JArray {
		return true
	}
	return false
}

//IsObject 类型检测
//IsObject test if is a object json element
func (node *JsonNode) IsObject() bool {
	if node.nodeType == JObject {
		return true
	}
	return false
}

//IsBool 类型检测
//IsBool test if is a bool json element
func (node *JsonNode) IsBool() bool {
	if node.nodeType == JBool {
		return true
	}
	return false
}

//IsNumber 类型检测
//IsNumber test if is a number json element
func (node *JsonNode) IsNumber() bool {
	if node.nodeType == JNumber {
		return true
	}
	return false
}

//IsString 类型检测
//IsString test if is a string json element
func (node *JsonNode) IsString() bool {
	if node.nodeType == JString {
		return true
	}
	return false
}

//IsNull 类型检测
//IsNull test if is a null json element
func (node *JsonNode) IsNull() bool {
	if node.nodeType == JNull {
		return true
	}
	return false
}

//Do 迭代
//Do iterater
func (node *JsonNode) Do(fn func(interface{}, *JsonNode)) error {
	switch node.nodeType {
	case JArray:
		for i, v := range node.value.(JArrayValue) {
			fn(i, v)
		}
	case JObject:
		for k, v := range node.value.(JObjectValue) {
			fn(k, v)
		}
	default:
		return errors.New("can not be iterated")
	}
	return nil
}

//Marshal 将当前json对象转换为json字符串
//Marshal make a json string
func (node *JsonNode) Marshal() (str string, err error) {
	commaFlag := false
	switch node.nodeType {
	case JArray:
		str = `[`
		for _, v := range node.value.(JArrayValue) {
			var tstr string
			tstr, err = v.Marshal()
			if err != nil {
				return "", err
			}
			if commaFlag {
				str = str + `,`
			}
			commaFlag = true
			str = str + tstr
		}
		str = str + `]`
	case JObject:
		str = `{`
		for k, v := range node.value.(JObjectValue) {
			var tstr string
			tstr, err = v.Marshal()
			if err != nil {
				return "", err
			}
			if commaFlag {
				str = str + `,`
			}
			commaFlag = true
			str = str + `"` + k + `":` + tstr
		}
		str = str + `}`
	case JString:
		str = `"` + node.value.(string) + `"`
	case JBool:
		if node.value.(bool) {
			str = "true"
		} else {
			str = "false"
		}
	case JNumber:
		str = strconv.FormatFloat(node.value.(float64), 'f', -1, 64)
	case JNull:
		str = "null"
	default:
		err = errors.New("wrong type")
	}
	return str, err
}
