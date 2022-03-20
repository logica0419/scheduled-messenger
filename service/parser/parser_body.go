package parser

import "strings"

// body から特定のルールにマッチする文字列を変換
func bodyParse(body *string) *string {
	// メンションよけのパース (@.{id} を @{id} に変換)
	replacedBody := strings.Replace(*body, "@.", "@", -1)

	return &replacedBody
}
