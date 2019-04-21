package game

// Interface for a generic clickable thing
// e.g. buttons/icons

var (
	// Keep track of whether there have ever been touches
	haveBeenTouches bool
)

type Clickable interface {
	Click() error
}

// The UI elements would call the "Click()" function on any clickable children
// e.g. to click a button in a window...
// root ui would call window.Click(), which would in turn call button.Click()
