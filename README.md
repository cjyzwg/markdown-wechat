# markdown-wechat
markdown文件转成微信公众号格式，根据自己的css添加

### Demo[链接](https://github.com/cjyzwg/markdown-wechat-demo)

#### **微信公众号配置项**
- cp .env.example .env
- 将.env配置修改成微信公众号对应的值
```
APPID=xxxxxxx
SECRET=xxxxxxxxxxxxx
TOKEN=xxxxx
EncodingAESKey=xxxxxxxxxxxxxxxxxxxx

```

#### **引用步骤**

```
import (
	"github.com/cjyzwg/markdown-wechat/draft"
	"github.com/fastwego/offiaccount"
	"github.com/spf13/viper"
)
```
```
res, err := draft.DraftRun(config_file, App)
```




