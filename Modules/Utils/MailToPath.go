package Utils

import "strings"

func MailToPath(mail string) string {
	var formatted = strings.ReplaceAll(mail, "@", "-")
	formatted = strings.ReplaceAll(formatted, ".", "-")
	formatted = strings.ReplaceAll(formatted, "_", "-")
	return formatted
}
