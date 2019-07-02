package log

var (
	// Black is the color black
	Black Color = "\033[00;30m"

	// Red terminal output color
	Red Color = "\033[00;31m"

	// Green terminal output color
	Green Color = "\033[00;32m"

	// Orange terminal output color
	Orange Color = "\033[00;33m"

	// Blue terminal output color
	Blue Color = "\033[00;34m"

	// Purple terminal output color
	Purple Color = "\033[00;35m"

	// Cyan terminal output color
	Cyan Color = "\033[00;36m"

	// LightGrey terminal output color
	LightGrey Color = "\033[00;37m"

	// DarkGrey terminal output color
	DarkGrey Color = "\033[01;30m"

	// LightRed terminal output color
	LightRed Color = "\033[01;31m"

	// LightGreen terminal output color
	LightGreen Color = "\033[01;32m"

	// Yellow terminal output color
	Yellow Color = "\033[01;33m"

	// LightBlue terminal output color
	LightBlue Color = "\033[01;34m"

	// LightPurple terminal output color
	LightPurple Color = "\033[01;35m"

	// LightCyan terminal output color
	LightCyan Color = "\033[01;36m"

	// White terminal output color
	White Color = "\033[01;37m"

	// NoColor will reset the current terminal color
	NoColor Color = "\033[0m"
)

// Color is a terminal color
type Color string

func (c Color) String() string {
	return string(c) // fmt.Sprintf("\033[%sm", c)
}
