package mail

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
	"io/ioutil"
	"sync"
	"time"
)

var once sync.Once

type Sender struct {
	Host        string // smtp.163.com
	Port        int    // 465
	User        string // 邮箱账号:youlink_yi
	Pwd         string // 邮箱发送密码
	ReTryTimes  int    // 邮件发送失败重试次数.默认10次
	ReTryPeriod int    // 邮件发送失败重试的间隔时间.默认60s
	dialer      *gomail.Dialer
	tmp         *gomail.Message
	ch          chan Message
	err         chan SenderError
}

type Message struct {

	// 发件人
	From User

	// 收件人列表
	To []User

	// 抄送
	Cc []User

	// 邮件标题
	Subject string

	// 附件
	Attach string

	// 邮件正文模板(html/template)
	Template string

	// 邮件正文模板数据对象
	TemplateObject interface{}
}

type User struct {
	Name string
	Addr string
}

type SenderError struct {
	Err error
	msg Message
}

func NewSender(host, user, pwd string, port int) *Sender {
	return &Sender{
		Host:        host,
		User:        user,
		Pwd:         pwd,
		Port:        port,
		ReTryTimes:  10,
		ReTryPeriod: 60,
		ch:          make(chan Message, 10),
		err:         make(chan SenderError, 10),
	}
}

// 发送邮件不需要大并发,所以使用单协程阻塞式发送
func (s *Sender) run() {
	go func(s *Sender) {
		var err error
		var t *template.Template

		for m := range s.ch {

			s.tmp = gomail.NewMessage()
			s.tmp.SetAddressHeader("From", m.From.Addr, m.From.Name)

			for _, v := range m.To {
				s.tmp.SetAddressHeader("To", v.Addr, v.Name)
			}

			for _, v := range m.Cc {
				s.tmp.SetAddressHeader("Cc", v.Addr, v.Name)
			}

			s.tmp.SetHeader("Subject", m.Subject)

			if m.Attach != "" {
				s.tmp.Attach(m.Attach)
			}

			t, err = template.New("mail").Parse(m.Template)
			if err != nil {
				s.err <- SenderError{Err: err, msg: m}
				continue
			}

			buf := new(bytes.Buffer)
			err = t.Execute(buf, m.TemplateObject)
			if err != nil {
				s.err <- SenderError{Err: err, msg: m}
				continue
			}
			var b []byte
			b, err = ioutil.ReadAll(buf)
			if err != nil {
				s.err <- SenderError{Err: err, msg: m}
				continue
			}

			s.tmp.SetBody("text/html", string(b))

			for i := 0; i < s.ReTryTimes; i++ {
				err = s.dialer.DialAndSend(s.tmp)
				if err != nil {
					s.err <- SenderError{Err: err, msg: m}
					time.Sleep(time.Duration(s.ReTryPeriod) * time.Second)
					continue
				}
				break
			}
		}
		close(s.err)
	}(s)
}

func (s *Sender) SendMail(m Message) {
	s.dialer = gomail.NewDialer(s.Host, s.Port, s.User, s.Pwd)
	once.Do(s.run)
	s.ch <- m
}

func (s *Sender) Err() <-chan SenderError {
	return s.err
}

func (s *Sender) Stop() {
	close(s.ch)
}
