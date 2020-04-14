package tray

type (
	TrayOption struct {
		Icon               string
		Title              string // not supported on Windows
		ToolTip            string
		Invisible          bool
		LeftClickCallback  func() // not supported on Linux
		RightClickCallback func() // not supported on Linux
		Children           []ItemOption
	}

	ItemOption struct {
		Icon      string
		Title     string
		Disabled  bool
		Checkable bool
		Checked   bool
		Invisible bool
		Callback  func() // On Linux, callback() when submenu is selected
		Children  []ItemOption
	}

	Tray interface {
		Run()
		SetIcon(icon string) error
		SetTitle(title string) error // not supported on Windows
		SetToolTip(toolTip string) error
		SetVisible(visible bool) error
		SetLeftClickCallback(callback func()) error  // not supported on Linux
		SetRightClickCallback(callback func()) error // not supported on Linux
		AddItem(option ItemOption) (Item, error)
		AddItems(options ...ItemOption) ([]Item, error)
		Close() error
	}

	Item interface {
		Checked() bool
		SetIcon(img string) error
		SetTitle(title string) error
		SetEnabled(enabled bool) error
		SetVisible(visible bool) error
		SetChecked(checked bool) error
		SetCallback(callback func()) error
		AddItem(option ItemOption) (Item, error)
		AddItems(options ...ItemOption) ([]Item, error)
	}
)

func NewTray(option TrayOption) (Tray, error) {
	return newTray(option)
}
