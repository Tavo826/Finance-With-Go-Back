package mail

import (
	"html/template"
	"log"
	"personal-finance/adapter/config"
	"personal-finance/core/domain"

	"github.com/wneessen/go-mail"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type MailReportAdapter struct {
	config *config.Mail
}

func NewMailReportAdapter(config *config.Mail) *MailReportAdapter {

	return &MailReportAdapter{
		config,
	}
}

func formatMoney(amount float64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("$%.2f", amount)
}

const (
	htmlBodyTemplate = `<p>&#10024; Hi {{.Username}} we hope you are having a great day!</p>

<p>Here is the monthly summary of your personal finances for <strong>{{.Month}}</strong> of <strong>{{.Year}}</strong>.</p>

<p>We invite you to continue recording all your expenses and incomes in the app.</p>

<p>&#128640; You can login in the following link: <a href="https://tavo826.github.io/Finance-With-Angular-Front/Home" target="_blank">Personal Finance</a></p>

<div style="display: flex; justify-content: center; margin-top: 20px;">
  <table style="border-collapse: collapse; font-family: Arial, sans-serif; width: 50%; box-shadow: 0 0 10px rgba(0,0,0,0.1);">
    <tr style="background-color: #f2f2f2;">
      <td style="border: 1px solid #ddd; padding: 12px; font-weight: bold;">Total income</td>
      <td style="border: 1px solid #ddd; padding: 12px;">{{formatMoney .TotalIncome}}</td>
    </tr>
    <tr>
      <td style="border: 1px solid #ddd; padding: 12px; font-weight: bold;">Total expenses</td>
      <td style="border: 1px solid #ddd; padding: 12px;">{{formatMoney .TotalExpenses}}</td>
    </tr>
    <tr style="background-color: #f9f9f9;">
      <td style="border: 1px solid #ddd; padding: 12px; font-weight: bold;">Net balance</td>
      <td style="border: 1px solid #ddd; padding: 12px;">{{formatMoney .NetBalance}}</td>
    </tr>
  </table>
</div>

<h3 style="text-align: center; margin-top: 40px;">Details by origin</h3>

<div style="display: flex; justify-content: center; margin-top: 10px;">
  <table style="border-collapse: collapse; font-family: Arial, sans-serif; width: 80%; box-shadow: 0 0 10px rgba(0,0,0,0.1);">
    <tr style="background-color: #f2f2f2;">
      <th style="border: 1px solid #ddd; padding: 12px;">Origin</th>
      <th style="border: 1px solid #ddd; padding: 12px;">Income</th>
      <th style="border: 1px solid #ddd; padding: 12px;">Expenses</th>
      <th style="border: 1px solid #ddd; padding: 12px;">Balance</th>
    </tr>
    {{range .OriginSummary}}
    <tr>
      <td style="border: 1px solid #ddd; padding: 12px;">{{.OriginName}}</td>
      <td style="border: 1px solid #ddd; padding: 12px;">{{formatMoney .TotalIncome}}</td>
      <td style="border: 1px solid #ddd; padding: 12px;">{{formatMoney .TotalExpenses}}</td>
      <td style="border: 1px solid #ddd; padding: 12px;">{{formatMoney .OriginBalance}}</td>
    </tr>
    {{end}}
  </table>
</div>`
)

func (ra *MailReportAdapter) SendMail(report domain.Report) error {

	message := mail.NewMsg()

	htmlTpl, err := template.New("htmltpl").Funcs(template.FuncMap{
		"formatMoney": formatMoney,
	}).Parse(htmlBodyTemplate)
	if err != nil {
		return err
	}

	if err := message.EnvelopeFrom(ra.config.Username); err != nil {
		return err
	}
	if err := message.FromFormat("Personal finance", ra.config.Username); err != nil {
		return err
	}
	if err := message.AddToFormat(report.Username, report.UserEmail); err != nil {
		log.Print("Error con el correo: ", err)
		return err
	}

	message.Subject("¡Toc-toc! Your personal finance summary is here")
	if err := message.AddAlternativeHTMLTemplate(htmlTpl, report); err != nil {
		return err
	}

	client, err := mail.NewClient(ra.config.Host, mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(ra.config.Username), mail.WithPassword(ra.config.Password))
	if err != nil {
		return err
	}

	if err := client.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
