package robot

import (
	"fmt"
	"strings"

	"github.com/opentdp/wechat-rest/dbase/keyword"
	"github.com/opentdp/wechat-rest/dbase/profile"
	"github.com/opentdp/wechat-rest/dbase/tables"
	"github.com/opentdp/wechat-rest/wcferry"
)

var badMember = map[string]int{}
var keywordList = []*tables.Keyword{}

func badHandler() {

	updateBadWord()

	handlers["/bad"] = &Handler{
		Level:    7,
		Order:    30,
		ChatAble: true,
		RoomAble: true,
		Describe: "添加违规关键词",
		Callback: func(msg *wcferry.WxMsg) string {
			_, err := keyword.Create(&keyword.CreateParam{
				Roomid: prid(msg), Phrase: msg.Content, Target: "ban", Level: 1,
			})
			if err == nil {
				updateBadWord()
				return "违规关键词添加成功"
			}
			return "违规关键词已存在"
		},
		PreCheck: badPreCheck,
	}

	handlers["/unbad"] = &Handler{
		Level:    7,
		Order:    31,
		ChatAble: true,
		RoomAble: true,
		Describe: "删除违规关键词",
		Callback: func(msg *wcferry.WxMsg) string {
			item, err := keyword.Fetch(&keyword.FetchParam{Roomid: prid(msg), Phrase: msg.Content})
			if err != nil || item.Target != "ban" {
				return msg.Content + "不是违规关键词"
			}
			err = keyword.Delete(&keyword.DeleteParam{Rd: item.Rd})
			if err == nil {
				updateBadWord()
				return "已删除违规关键词" + msg.Content
			}
			return "违规关键词删除失败"
		},
	}

}

func badPreCheck(msg *wcferry.WxMsg) string {

	// 私聊豁免
	if !msg.IsGroup {
		return ""
	}

	// 管理豁免
	up, _ := profile.Fetch(&profile.FetchParam{Wxid: msg.Sender, Roomid: prid(msg)})
	if up.Level >= 7 {
		return ""
	}

	// 遍历关键词
	for _, v := range keywordList {
		if v.Roomid != "-" {
			if msg.IsGroup {
				if v.Roomid != msg.Roomid {
					continue
				}
			} else {
				continue
			}
		}
		if v.Level > 0 && strings.Contains(msg.Content, v.Phrase) {
			badMember[msg.Sender] += int(v.Level)
			if badMember[msg.Sender] > 10 {
				wc.CmdClient.DelChatRoomMembers(msg.Roomid, msg.Sender)
				delete(badMember, msg.Sender)
				return "我送你离开，天涯之外你是否还在"
			}
			str := "违规风险 +%d，当前累计：%d，大于 10 将被请出群聊"
			return fmt.Sprintf(str, v.Level, badMember[msg.Sender])
		}
	}

	return ""

}

func updateBadWord() {

	list, _ := keyword.FetchAll(&keyword.FetchAllParam{})
	keywordList = list

}
