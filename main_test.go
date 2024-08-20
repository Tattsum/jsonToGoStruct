package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateStruct(t *testing.T) {
	testCases := []struct {
		name       string
		structName string
		jsonObject map[string]interface{}
		expected   Struct[string]
	}{
		{
			name:       "基本的な構造体生成",
			structName: "TestStruct",
			jsonObject: map[string]interface{}{
				"id":   1,
				"name": "テスト",
				"age":  30,
			},
			expected: Struct[string]{
				Name: "TestStruct",
				Fields: []Field[string]{
					{Name: "Id", Type: "float64", Tag: "`json:\"id\"`"},
					{Name: "Name", Type: "string", Tag: "`json:\"name\"`"},
					{Name: "Age", Type: "float64", Tag: "`json:\"age\"`"},
				},
			},
		},
		{
			name:       "ネストされたオブジェクトを含む構造体生成",
			structName: "NestedStruct",
			jsonObject: map[string]interface{}{
				"user": map[string]interface{}{
					"name":  "ユーザー1",
					"email": "user1@example.com",
				},
				"is_active": true,
			},
			expected: Struct[string]{
				Name: "NestedStruct",
				Fields: []Field[string]{
					{Name: "User", Type: "User", Tag: "`json:\"user\"`"},
					{Name: "IsActive", Type: "bool", Tag: "`json:\"is_active\"`"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GenerateStruct(tc.structName, tc.jsonObject)
			assert.Equal(t, tc.expected.Name, result.Name, "構造体名が一致しません")
			assert.Equal(t, len(tc.expected.Fields), len(result.Fields), "フィールド数が一致しません")
			for i, expectedField := range tc.expected.Fields {
				assert.Equal(t, expectedField.Name, result.Fields[i].Name, "フィールド名が一致しません")
				assert.Equal(t, expectedField.Type, result.Fields[i].Type, "フィールドの型が一致しません")
				assert.Equal(t, expectedField.Tag, result.Fields[i].Tag, "フィールドのタグが一致しません")
			}
		})
	}
}
