package tray

import (
	"runtime"

	"github.com/dawidd6/go-appindicator"
	"github.com/gotk3/gotk3/gtk"
)

type linuxTray struct {
	indicator *appindicator.Indicator
	mainMenu  linuxItem
}

func newTray(option TrayOption) (Tray, error) {

	runtime.LockOSThread()
	gtk.Init(nil)

	var (
		res linuxTray
		err error
	)

	res.indicator = appindicator.New(option.Title, "", appindicator.CategoryApplicationStatus)

	res.mainMenu.menu, err = gtk.MenuNew()
	if err != nil {
		return nil, err
	}
	res.indicator.SetMenu(res.mainMenu.menu)

	if option.Icon != "" {
		if err = res.SetIcon(option.Icon); err != nil {
			return nil, err
		}
	}

	if option.Title != "" {
		if err = res.SetTitle(option.Title); err != nil {
			return nil, err
		}
	}

	if option.ToolTip != "" {
		if err = res.SetToolTip(option.ToolTip); err != nil {
			return nil, err
		}
	}

	if !option.Invisible {
		if err = res.SetVisible(!option.Invisible); err != nil {
			return nil, err
		}
	}

	if option.LeftClickCallback != nil {
		if err = res.SetLeftClickCallback(option.LeftClickCallback); err != nil {
			return nil, err
		}
	}

	if option.RightClickCallback != nil {
		if err = res.SetRightClickCallback(option.LeftClickCallback); err != nil {
			return nil, err
		}
	}

	if len(option.Children) > 0 {
		_, err = res.AddItems(option.Children...)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (t linuxTray) SetIcon(icon string) error {
	iconPath, err := getIconPath(icon)
	if err != nil {
		return err
	}
	t.indicator.SetIconFull(iconPath, "")
	return nil
}

func (t linuxTray) Run() {
	runtime.LockOSThread()
	t.mainMenu.menu.ShowAll()
	gtk.Main()
}

func (t linuxTray) SetTitle(title string) error {
	t.indicator.SetTitle(title)
	return nil
}

func (t linuxTray) SetToolTip(toolTip string) error {
	t.indicator.SetLabel(toolTip, "")
	return nil
}

func (t linuxTray) SetVisible(visible bool) error {
	if visible {
		t.indicator.SetStatus(appindicator.StatusActive)
	} else {
		t.indicator.SetStatus(appindicator.StatusPassive)
	}
	return nil
}

func (t linuxTray) SetLeftClickCallback(func()) error {
	return nil
}

func (t linuxTray) SetRightClickCallback(func()) error {
	return nil
}

func (t linuxTray) AddItem(option ItemOption) (Item, error) {
	return t.mainMenu.AddItem(option)
}

func (t linuxTray) AddItems(options ...ItemOption) ([]Item, error) {
	return t.mainMenu.AddItems(options...)
}

func (t linuxTray) Close() error {
	gtk.MainQuit()
	return nil
}
