```go
package main

import (
	"fmt"
	"github.com/pyting/sendmail"
	"time"
)

type Date struct {
	Title       string  // 基金标题
	Trader      string  // 操盘手姓名
	Total       float64 // 汇总金额
	SubUser     string  // 火币子帐号ID
	SubUserName string  // 火币子帐号名称
}

func main() {

	// 每日基金购入汇总模板
	t1 := `
<table frame=box rules=rows bgcolor=#FFFAFA cellpadding=20 cellspacing=10>
<tr>
<th>编号</th><th>基金标题</th><th>操盘手</th><th>汇总金额</th><th>火币子帐号ID</th><th>火币子帐号名称</th>
</tr>
{{range $i,$v := .}}
<tr>
<th>{{$i}}</th><th>{{$v.Title}}</th><th>{{$v.Trader}}</th><th>{{$v.Total}}</th><th>{{$v.SubUser}}</th><th>{{$v.SubUserName}}</th>
</tr>
{{end}}
</table>
`

	sender := mail.NewSender("smtp.163.com", "youlink_yi", "Yml4dWFuCg", 465)

	from := mail.User{"Bixuan", "youlink_yi@163.com"}
	to := []mail.User{{"胡一晟", "2312493805@qq.com"}}
	cc := []mail.User{{"张一斌", "382431937@qq.com"}}

	data := []Date{
		{"币种多样，收益稳健", "胡伟", 19874.3454, "102528316", "tumen003"},
		{"币种多样，收益稳健", "胡伟", 19874.3454, "102528316", "tumen003"},
		{"币种多样，收益稳健", "胡伟", 19874.3454, "102528316", "tumen003"},
		{"币种多样，收益稳健", "胡伟", 19874.3454, "102528316", "tumen003"},
		{"币种多样，收益稳健", "胡伟", 19874.3454, "102528316", "tumen003"},
		{"币种多样，收益稳健", "胡伟", 19874.3454, "102528316", "tumen003"},
		{"币种多样，收益稳健", "胡伟", 19874.3454, "102528316", "tumen003"},
		{"币种多样，收益稳健", "胡伟", 19874.3454, "102528316", "tumen003"},
		{"币种多样，收益稳健", "胡伟", 19874.3454, "102528316", "tumen003"},
		{"币种多样，收益稳健", "胡伟", 19874.3454, "102528316", "tumen003"},
		{"币种多样，收益稳健", "胡伟", 19874.3454, "102528316", "tumen003"},
		{"币种多样，收益稳健", "胡伟", 19874.3454, "102528316", "tumen003"},
	}

	m := mail.Message{
		From:           from,
		To:             to,
		Cc:             cc,
		Subject:        time.Now().Format("2006-01-02 15:04:05") + " 每日基金购入汇总",
		Template:       t1,
		TemplateObject: data,
	}

	sender.SendMail(m)

	// 异步获取错误信息
	go func() {
		for err := range sender.Err() {
			fmt.Println(err)
		}
	}()

	time.Sleep(10 * time.Second)

	sender.Stop()
}

```