package draft

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/fastwego/offiaccount"
	"github.com/fastwego/offiaccount/apis/material"
)

func DraftRun(config_file *ConfigFile, App *offiaccount.OffiAccount) (string, error) {

	content := MarkdownRun(config_file.MarkdownFilePath, config_file.CssFilePath, App)

	//修改空格的问题
	// content = strings.Replace(content, "@s-;", "&nbsp;", -1)

	content, _ = HtmlMinify(content)

	// 	//新增图片素材，获取media_id
	params := url.Values{}
	params.Add("type", "thumb")
	fields := map[string]string{}
	resp, err := material.AddMaterial(App, config_file.ImagePath, params, fields)
	fmt.Println(string(resp), err)

	var thumb_return ThumbReturn
	if err := json.Unmarshal([]byte(resp), &thumb_return); err == nil {
		fmt.Println(thumb_return)
	} else {
		// fmt.Println(err)
		// return "永久图文素材上传错误"
		return "", err
	}

	articles := make(map[string][]Article)
	article := Article{
		Title:              config_file.Title,
		Author:             config_file.Author,
		Digest:             config_file.Digest,
		Content:            content,
		ContentSourceUrl:   "",
		ThumbMediaId:       thumb_return.MediaId,
		NeedOpenComment:    0,
		OnlyFansCanComment: 0,
	}
	var json_articles []Article
	json_articles = append(json_articles, article)
	articles["articles"] = json_articles

	payload, err := json.Marshal(articles)
	if err != nil {
		// fmt.Println(err)
		return "", err
	}
	resp, err = AddDraft(App, payload)
	// fmt.Println(string(resp), err)
	return string(resp), err
}
