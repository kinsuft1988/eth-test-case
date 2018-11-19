package blocReport

import "gopkg.in/gomail.v2"

type Report struct {
	address string
}

func (* Report)SendMail(msg string)  {
	m := gomail.NewMessage()
	m.SetHeader("From", "475002739@qq.com")
	m.SetHeader("To", "khu@block-cloud.cn")
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Eth 性能测试")
	m.SetBody("text/html", msg)
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.qq.com", 465, "475002739", "afzbtsshqxbpbhbd")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}