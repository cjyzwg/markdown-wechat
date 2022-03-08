// Package draft 草稿箱的api
package draft

import (
	"bytes"

	"github.com/fastwego/offiaccount"
)

const (
	apiAddDraft      = "/cgi-bin/draft/add"
	apiGetDraft      = "/cgi-bin/draft/get"
	apiDelDraft      = "/cgi-bin/draft/delete"
	apiUpdateDraft   = "/cgi-bin/draft/update"
	apiGetDraftCount = "/cgi-bin/draft/count"
	apiBatchgetDraft = "/cgi-bin/draft/batchget"
)

/*
新增草稿箱

新增常用的素材到草稿箱中进行使用。上传到草稿箱中的素材被群发或发布后，该素材将从草稿箱中移除。新增草稿可在公众平台官网-草稿箱中查看和管理

See: https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Add_draft.html

POST https://api.weixin.qq.com/cgi-bin/draft/add?access_token=ACCESS_TOKEN
*/
func AddDraft(ctx *offiaccount.OffiAccount, payload []byte) (resp []byte, err error) {
	return ctx.Client.HTTPPost(apiAddDraft, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
获取草稿

新增草稿后，开发者可以根据草稿指定的字段来下载草稿。


See: https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Get_draft.html

POST https://api.weixin.qq.com/cgi-bin/draft/get?access_token=ACCESS_TOKEN
*/
func GetDraft(ctx *offiaccount.OffiAccount, payload []byte) (resp []byte, err error) {
	return ctx.Client.HTTPPost(apiGetDraft, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
删除草稿

新增草稿后，开发者可以根据本接口来删除不再需要的草稿，节省空间。此操作无法撤销，请谨慎操作。

See: https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Delete_draft.html

POST https://api.weixin.qq.com/cgi-bin/draft/delete?access_token=ACCESS_TOKEN
*/
func DelDraft(ctx *offiaccount.OffiAccount, payload []byte) (resp []byte, err error) {
	return ctx.Client.HTTPPost(apiDelDraft, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
修改草稿

开发者可通过本接口对草稿进行修改。


See: https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Update_draft.html

POST https://api.weixin.qq.com/cgi-bin/draft/update?access_token=ACCESS_TOKEN
*/
func UpdateDraf(ctx *offiaccount.OffiAccount, payload []byte) (resp []byte, err error) {
	return ctx.Client.HTTPPost(apiUpdateDraft, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
获取草稿总数

开发者可以根据本接口来获取草稿的总数。此接口只统计数量，不返回草稿的具体内容。


See: https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Count_drafts.html

GET https://api.weixin.qq.com/cgi-bin/draft/count?access_token=ACCESS_TOKEN

*/
func GetDraftCount(ctx *offiaccount.OffiAccount) (resp []byte, err error) {
	return ctx.Client.HTTPGet(apiGetDraftCount)
}

/*
获取草稿列表

新增草稿之后，开发者可以获取草稿的列表。

See: https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Get_draft_list.html

POST https://api.weixin.qq.com/cgi-bin/draft/batchget?access_token=ACCESS_TOKEN
*/
func BatchgetDraft(ctx *offiaccount.OffiAccount, payload []byte) (resp []byte, err error) {
	return ctx.Client.HTTPPost(apiBatchgetDraft, bytes.NewReader(payload), "application/json;charset=utf-8")
}
