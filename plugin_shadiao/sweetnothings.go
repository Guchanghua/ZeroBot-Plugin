package shadiao

import (
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	"github.com/FloatTech/zbputils/web"
)

func init() {
	engine.OnFullMatch("来碗绿茶").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if !limit.Load(ctx.Event.GroupID).Acquire() {
			return
		}
		data, err := web.ReqWith(chayiURL, "GET", loveliveReferer, ua)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		text := gjson.Get(helper.BytesToString(data), "returnObj.content").String()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(text))
	})

	engine.OnFullMatch("渣我").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if !limit.Load(ctx.Event.GroupID).Acquire() {
			return
		}
		data, err := web.ReqWith(ganhaiURL, "GET", loveliveReferer, ua)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		text := gjson.Get(helper.BytesToString(data), "returnObj.content").String()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(text))
	})
}
