package tray

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type commonItem interface {
	Add(gtk.IWidget)
	SetSubmenu(gtk.IWidget)
	gtk.IWidget
	Show()
	SetNoShowAll(bool)
	Hide()
	SetVisible(bool)
	SetSensitive(bool)
	Connect(detailedSignal string, f interface{}, userData ...interface{}) (glib.SignalHandle, error)
	HandlerDisconnect(handle glib.SignalHandle)
}

type linuxItem struct {
	menu     *gtk.Menu
	delegate commonItem
	box      *gtk.Box
	label    *gtk.Label
	icon     *gtk.Image
	handle   *glib.SignalHandle
}

func (i *linuxItem) Checked() bool {
	if checkitem, ok := i.delegate.(*gtk.CheckMenuItem); ok {
		return checkitem.GetActive()
	}
	return false
}

func (i *linuxItem) SetIcon(icon string) error {
	if i.icon == nil {
		img, err := gtk.ImageNew()
		if err != nil {
			return err
		}
		if i.label != nil {
			i.box.Remove(i.label)
		}
		i.icon = img
		i.box.Add(i.icon)
		i.box.Add(i.label)
	}
	if icon != "-" {
		iconPath, err := getIconPath(icon)
		if err != nil {
			return err
		}
		i.icon.SetFromFile(iconPath)
	}
	return nil
}

func (i *linuxItem) SetTitle(title string) error {
	if i.label == nil {
		label, err := gtk.LabelNew(title)
		if err != nil {
			return err
		}
		i.label = label
		i.box.Add(i.label)
		return nil
	}
	i.label.SetLabel(title)
	return nil
}

func (i *linuxItem) SetEnabled(enabled bool) error {
	i.delegate.SetSensitive(enabled)
	return nil
}

func (i *linuxItem) SetVisible(visible bool) error {
	i.delegate.SetNoShowAll(!visible)
	return nil
}

func (i *linuxItem) SetChecked(checked bool) error {
	if checkitem, ok := i.delegate.(*gtk.CheckMenuItem); ok {
		checkitem.SetActive(checked)
	}
	return nil
}

func (i *linuxItem) SetCallback(callback func()) error {
	if i.handle != nil {
		i.delegate.HandlerDisconnect(*i.handle)
	}
	if callback != nil {
		signalName := "activate"
		if _, ok := i.delegate.(*gtk.CheckMenuItem); ok {
			signalName = "toggled"
		}
		handle, err := i.delegate.Connect(signalName, callback)
		if err != nil {
			return err
		}
		i.handle = &handle
	}
	return nil
}

func (i *linuxItem) AddItem(option ItemOption) (Item, error) {
	if i.menu == nil {
		menu, err := gtk.MenuNew()
		if err != nil {
			return nil, err
		}
		i.delegate.(*gtk.MenuItem).SetSubmenu(menu)
		i.menu = menu
	}
	item, err := newItem(option)
	if err != nil {
		return nil, err
	}
	i.menu.Add(item.delegate)
	return item, nil
}

func (i *linuxItem) AddItems(options ...ItemOption) ([]Item, error) {
	resLen := len(options)
	if resLen == 0 {
		return nil, nil
	}
	res := make([]Item, 0, resLen)
	for _, v := range options {
		item, err := i.AddItem(v)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func newItem(option ItemOption) (*linuxItem, error) {
	var (
		res linuxItem
		err error
	)
	if option.Title == "" && option.Icon == "" {
		res.delegate, err = gtk.SeparatorMenuItemNew()
		return &res, err
	}

	var item commonItem
	if option.Checkable {
		item, err = gtk.CheckMenuItemNew()
		if err != nil {
			return nil, err
		}
	} else {
		item, err = gtk.MenuItemNew()
		if err != nil {
			return nil, err
		}
	}

	res.delegate = item
	if option.Checked {
		if err = res.SetChecked(option.Checked); err != nil {
			return nil, err
		}
	}

	res.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		return nil, err
	}
	item.Add(res.box)

	if option.Children != nil {
		res.menu, err = gtk.MenuNew()
		if err != nil {
			return nil, err
		}
		item.SetSubmenu(res.menu)

		_, err = res.AddItems(option.Children...)
		if err != nil {
			return nil, err
		}
	}

	if option.Title != "" {
		if err = res.SetTitle(option.Title); err != nil {
			return nil, err
		}
	}

	if option.Disabled {
		if err = res.SetEnabled(false); err != nil {
			return nil, err
		}
	}

	if option.Invisible {
		if err = res.SetVisible(false); err != nil {
			return nil, err
		}
	}

	if option.Icon != "" {
		if err := res.SetIcon(option.Icon); err != nil {
			return nil, err
		}
	}

	if option.Callback != nil {
		if err := res.SetCallback(option.Callback); err != nil {
			return nil, err
		}
	}

	return &res, nil
}
