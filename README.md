####Examples
#####code:
```go
func main() {
	s := `{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`
	node, _ := jsonhandler.NewJsonNode()
	err := node.Unmarshal([]byte(s))
	a, err := node.GetString("Parents", 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)
}
```
#####output:
```
Gomez
```

#####code:
```go
func main() {
	node, err := jsonhandler.NewJsonNode()
	node.Set(1, "A")
	node.Set("b", "B")
	a, err := node.Marshal()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)
}
```
#####output:
```
{"A":1,"B":"b"}
```

#####code:
```go
func main() {
	s := `[{"Name":"adam","Age":6},{"Name":"eve","Age":7}]`
	node, err := jsonhandler.NewJsonNode()
	err = node.Unmarshal([]byte(s))
	node.Set(10, 1, "age")
	a, err := node.GetNumber(1, "Age")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)
	fmt.Println(node.Marshal())
}
```
#####output:
```
10
[{"Name":"adam","Age":6},{"Name":"eve","Age":10}] <nil>
```

#####code:
```go
func main() {
	s := `["hello","world","go"]`
	node, err := jsonhandler.NewJsonNode()
	_ = node.Unmarshal([]byte(s))
	node.Do(func(i interface{}, v *jsonhandler.JsonNode) {
		a1, _ := v.String()
		fmt.Println(a1)
	})
	if err != nil {
		fmt.Println(err)
	}
}
```
#####output:
```
hello
world
go
```

#####code:
```go
func main() {
	s1 := `{"Name":"adam","Age":6}`
	s2 := `{"Name":"eve","Age":7}`
	node1, err := jsonhandler.NewJsonNode()
	_ = node1.Unmarshal([]byte(s1))
	node2, err := jsonhandler.NewJsonNode()
	_ = node2.Unmarshal([]byte(s2))
	node, err := jsonhandler.NewJsonNode()
	_, err = node.Set(node1, 0)
	_, err = node.Set(node2, 1)
	if err != nil {
		fmt.Println(err)
	}
	a, err := node.Marshal()
	fmt.Println(a)
}
```
#####output:
```
[{"Name":"adam","Age":6},{"Age":7,"Name":"eve"}]
```