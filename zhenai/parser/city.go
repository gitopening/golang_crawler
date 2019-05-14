package parser

import (
	"crawler/engine"
	"regexp"
)

var cityRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`)
var sexRe = regexp.MustCompile(`<td width="180"><span class="grayL">性别：</span>([^<]+)</td>`)
var cityUrlRe = regexp.MustCompile(
	`href="(http://www.zhenai.com/zhenghun/[^"]+)"`)

// 城市页面用户解析器
func ParseCity(bytes []byte) engine.ParseResult {
	submatch := cityRe.FindAllSubmatch(bytes, -1)
	gendermatch := sexRe.FindAllSubmatch(bytes, -1)

	result := engine.ParseResult{}

	for k, item := range submatch {
		name := string(item[2])
		gender := string(gendermatch[k][1])

		result.Items = append(result.Items, "User:"+name)
		result.Requests = append(result.Requests, engine.Request{
			Url: string(item[1]),
			ParseFunc: func(bytes []byte) engine.ParseResult {
				return ParseProfile(bytes, name, gender)
			},
		})
	}

	// 添加更多城市
	matches := cityUrlRe.FindAllSubmatch(bytes, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests, engine.Request{
			Url:       string(m[1]),
			ParseFunc: ParseCity,
		})
	}
	return result
}
