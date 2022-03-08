package draft

type ThumbReturn struct {
	MediaId string        `json:"media_id"`
	Url     string        `json:"url"`
	Item    []interface{} `json:"item"`
}
type Articles map[string][]Article

type Article struct {
	Title              string `json:"title"`
	Author             string `json:"author"`
	Digest             string `json:"digest"`
	Content            string `json:"content"`
	ContentSourceUrl   string `json:"content_source_url"`
	ThumbMediaId       string `json:"thumb_media_id"`
	NeedOpenComment    int    `json:"need_open_comment"`
	OnlyFansCanComment int    `json:"only_fans_can_comment"`
}
