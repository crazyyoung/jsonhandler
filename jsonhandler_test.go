package jsonhandler

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestGetJsonObject(t *testing.T) {
	s := `{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`
	node, _ := NewJsonNode()
	_ = node.Unmarshal([]byte(s))
	a1, _ := node.GetString("Parents", 0)
	assert.Equal(t, a1, "Gomez")
}
func TestGetJsonArray(t *testing.T) {
	s := `[{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]},"hello"]`
	node, _ := NewJsonNode()
	_ = node.Unmarshal([]byte(s))
	a1, _ := node.GetString(0, "Parents", 0)
	assert.Equal(t, a1, "Gomez")
}
func TestMakeJsonObject(t *testing.T) {
	node, _ := NewJsonNode()
	node.Set(1, "A")
	node.Set("b", "B")
	a1, _ := node.Get("A")
	a11, _ := a1.Number()
	assert.Equal(t, a11, 1.0)
	a2, _ := node.Get("B")
	a21, _ := a2.String()
	assert.Equal(t, a21, "b")
	assert.Equal(t, node.IsObject(), true)
}
func TestMakeJsonArray(t *testing.T) {
	node, _ := NewJsonNode()
	node.Set(1, 0)
	node.Set("b", 1)
	a1, _ := node.GetNumber(0)
	assert.Equal(t, a1, 1.0)
	a2, _ := node.GetString(1)
	assert.Equal(t, a2, "b")
	assert.Equal(t, node.IsArray(), true)
}
func TestSetJsonObject(t *testing.T) {
	s := `{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`
	node, _ := NewJsonNode()
	_ = node.Unmarshal([]byte(s))
	node.Set("23", "Age")
	a1, _ := node.GetString("Age")
	assert.Equal(t, a1, "23")
	node.Set(23, "Parents", 4)
	a2, _ := node.GetNumber("Parents", 4)
	assert.Equal(t, a2, 23.0)
	a3, e3 := node.GetString("Parents", 4)
	assert.Equal(t, a3, "")
	assert.NotEqual(t, e3, nil)
	a4, e4 := node.GetNull("Parents", 3)
	assert.Equal(t, a4, nil)
	assert.Equal(t, e4, nil)
	s2 := []byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`)
	node2, _ := NewJsonNode()
	_ = node2.Unmarshal([]byte(s2))
	node.Set(node2, "Other")
	a5, _ := node.GetString("Other", "Body")
	assert.Equal(t, a5, "Hello")
	node.Set(true, "Other", "Body")
	a6, _ := node.GetBool("Other", "Body")
	assert.Equal(t, a6, true)
}
func TestSetJsonArray(t *testing.T) {
	s := `[{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]},"hello"]`
	node, _ := NewJsonNode()
	_ = node.Unmarshal([]byte(s))
	node.Set(2333, 3)
	a1, _ := node.GetNumber(3)
	assert.Equal(t, a1, 2333.0)
}
func TestDelJsonObject(t *testing.T) {
	s := `{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`
	node, _ := NewJsonNode()
	_ = node.Unmarshal([]byte(s))
	node.Delete(2, 4, 5)
	node.Delete("Parents", 0)
	a1, _ := node.GetString("Parents", 0)
	assert.Equal(t, a1, "Morticia")
	node.Delete("Name")
	node.Delete("Parents")
	a2, _ := node.Marshal()
	assert.Equal(t, a2, `{"Age":6}`)
}
func TestDelJsonArray(t *testing.T) {
	s := `[{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]},"hello"]`
	node, _ := NewJsonNode()
	_ = node.Unmarshal([]byte(s))
	node.Delete("2", "4", "5")
	node.Delete(0)
	a1, _ := node.GetString(0)
	assert.Equal(t, a1, "hello")
}

func TestIteratorObject(t *testing.T) {
	s := `{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`
	node, _ := NewJsonNode()
	_ = node.Unmarshal([]byte(s))
	node.Do(func(i interface{}, v *JsonNode) {
		if i.(string) == "Age" {
			a1, _ := v.Number()
			assert.Equal(t, a1, 6.0)
			v.Set(12)
		}
	})
	a2, _ := node.GetNumber("Age")
	assert.Equal(t, a2, 12.0)
}

func TestIteratorArray(t *testing.T) {
	s := `["hello","world","go"]`
	node, _ := NewJsonNode()
	_ = node.Unmarshal([]byte(s))
	node.Do(func(i interface{}, v *JsonNode) {
		if i.(int) == 1 {
			a1, _ := v.String()
			assert.Equal(t, a1, "world")
		}
	})
}

func TestMake(t *testing.T) {
	node1, _ := NewJsonNode()
	node1.Set("adam", "name")
	node1.Set(11, "age")
	node2, _ := NewJsonNode()
	node2.Set("eve", "name")
	node2.Set(12, "age")
	node12, _ := NewJsonNode()
	node12.Set(node1, 0)
	node12.Set(node2, 1)

	a1, _ := node12.GetString(0, "name")
	assert.Equal(t, a1, "adam")
	a2, _ := node12.GetNumber(1, "age")
	assert.Equal(t, a2, 12.0)
}
