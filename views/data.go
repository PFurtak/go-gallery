package views

import "github.com/Users/patrickfurtak/desktop/go-gallery/models"

const (
	// AlerLvls change the displayed alert's color
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlSuccess = "success"
	AlertLvlInfo    = "primary"

	// AlertTypes modify the initial message in the strong tag
	AlertTypeError   = "Error: "
	AlertTypeWarning = "Warning!"
	AlertTypeSuccess = "Success!"
	AlertTypeInfo    = "Info: "

	// AlertMessages modify the text displayed in the alert
	AlertMessageGeneric = "Something went wrong, please contact us if issue persists."
)

// Alert is used to render bootstrap alert messages
type Alert struct {
	Level     string
	AlertType string
	Message   string
}

// Data is the top level structure that views expects data to come in as
type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.Alert = &Alert{
			Level:     AlertLvlError,
			AlertType: AlertTypeError,
			Message:   pErr.Public(),
		}
	} else {
		d.Alert = &Alert{
			Level:     AlertLvlError,
			AlertType: AlertTypeError,
			Message:   AlertMessageGeneric,
		}
	}
}

type PublicError interface {
	error
	Public() string
}
