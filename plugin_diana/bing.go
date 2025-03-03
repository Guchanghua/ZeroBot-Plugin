// Package diana 虚拟偶像女团 A-SOUL 成员嘉然相关
package diana

import (
	"math/rand"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	control "github.com/FloatTech/zbputils/control"

	"github.com/FloatTech/ZeroBot-Plugin/order"
	"github.com/FloatTech/ZeroBot-Plugin/plugin_diana/data"
)

var engine = control.Register("diana", order.PrioDiana, &control.Options{
	DisableOnDefault: false,
	Help: "嘉然\n" +
		"- 小作文\n" +
		"- 发大病\n" +
		"- 教你一篇小作文[作文]\n" +
		"- [回复]查重",
})

func init() {
	// 随机发送一篇上面的小作文
	engine.OnFullMatch("小作文").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			rand.Seed(time.Now().UnixNano())
			// 绕过第一行发病
			ctx.SendChain(message.Text(data.RandText()))
		})
	// 逆天
	engine.OnFullMatch("发大病").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			// 第一行是发病
			ctx.SendChain(message.Text(data.HentaiText()))
		})
	// 增加小作文
	engine.OnRegex(`^教你一篇小作文(.*)$`, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			err := data.AddText(ctx.State["regex_matched"].([]string)[1])
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
			} else {
				ctx.SendChain(message.Text("记住啦!"))
			}
		})
}
