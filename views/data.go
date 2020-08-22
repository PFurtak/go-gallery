package views

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
	Yield interface{}
}
