package tray

import (
	"github.com/lxn/walk"
)

type windowsTray struct {
	mainWindow         *walk.MainWindow
	notifyIcon         *walk.NotifyIcon
	leftClickedHandle  *int
	rightClickedHandle *int
}

func newTray(option TrayOption) (Tray, error) {

	var (
		res windowsTray
		err error
	)

	if res.mainWindow, err = walk.NewMainWindow(); err != nil {
		return nil, err
	}

	if res.notifyIcon, err = walk.NewNotifyIcon(res.mainWindow); err != nil {
		return nil, err
	}

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
		if err = res.SetRightClickCallback(option.RightClickCallback); err != nil {
			return nil, err
		}
	}

	if option.Children != nil {
		_, err = res.AddItems(option.Children...)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (t windowsTray) SetIcon(icon string) error {
	iconPath, err := getIconPath(icon)
	if err != nil {
		return err
	}
	iconFile, err := walk.NewIconFromFile(iconPath)
	if err != nil {
		return err
	}
	return t.notifyIcon.SetIcon(iconFile)
}

func (t windowsTray) Run() {
	t.mainWindow.Run()
}

func (t windowsTray) SetTitle(string) error {
	return nil
}

func (t windowsTray) SetToolTip(toolTip string) error {
	return t.notifyIcon.SetToolTip(toolTip)
}

func (t windowsTray) SetVisible(visible bool) error {
	return t.notifyIcon.SetVisible(visible)
}

func (t windowsTray) SetLeftClickCallback(callback func()) error {
	if t.leftClickedHandle != nil {
		t.notifyIcon.MouseUp().Detach(*t.leftClickedHandle)
	}
	if callback != nil {
		t.notifyIcon.MouseUp().Attach(func(x, y int, button walk.MouseButton) {
			if button == walk.LeftButton {
				callback()
			}
		})
	}
	return nil
}

func (t windowsTray) SetRightClickCallback(callback func()) error {
	if t.rightClickedHandle != nil {
		t.notifyIcon.MouseUp().Detach(*t.rightClickedHandle)
	}
	if callback != nil {
		t.notifyIcon.MouseUp().Attach(func(x, y int, button walk.MouseButton) {
			if button == walk.RightButton {
				callback()
			}
		})
	}
	return nil
}

func (t windowsTray) AddItem(option ItemOption) (Item, error) {
	item, err := newItem(option)
	if err != nil {
		return nil, err
	}
	err = t.notifyIcon.ContextMenu().Actions().Add(item.action)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (t windowsTray) AddItems(options ...ItemOption) ([]Item, error) {
	resLen := len(options)
	if resLen == 0 {
		return nil, nil
	}
	res := make([]Item, 0, resLen)
	for _, v := range options {
		item, err := t.AddItem(v)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func (t windowsTray) Close() error {
	err := t.notifyIcon.Dispose()
	if err != nil {
		return err
	}
	return t.mainWindow.Close()
}
