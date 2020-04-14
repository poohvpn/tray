package tray

import (
	"errors"

	"github.com/lxn/walk"
)

type windowsItem struct {
	action        *walk.Action
	menu          *walk.Menu
	clickedHandle *int
}

func (i windowsItem) Checked() bool {
	return i.action.Checked()
}

func (i windowsItem) SetIcon(icon string) error {
	if icon == "-" {
		return nil
	}
	iconPath, err := getIconPath(icon)
	if err != nil {
		return err
	}
	iconFile, err := walk.NewIconFromFile(iconPath)
	if err != nil {
		return err
	}
	return i.action.SetImage(iconFile)
}

func (i windowsItem) SetTitle(title string) error {
	return i.action.SetText(title)
}

func (i windowsItem) SetEnabled(enabled bool) error {
	return i.action.SetEnabled(enabled)
}

func (i windowsItem) SetVisible(visible bool) error {
	return i.action.SetVisible(visible)
}

func (i windowsItem) SetCheckable(checkable bool) error {
	return i.action.SetCheckable(checkable)
}

func (i windowsItem) SetChecked(checked bool) error {
	return i.action.SetChecked(checked)
}

func (i windowsItem) SetCallback(callback func()) error {
	if i.clickedHandle != nil {
		i.action.Triggered().Detach(*i.clickedHandle)
	}
	if callback != nil {
		handle := i.action.Triggered().Attach(callback)
		i.clickedHandle = &handle
	}
	return nil
}

func (i windowsItem) AddItem(option ItemOption) (Item, error) {
	if i.menu == nil {
		return nil, errors.New("only menu can add item")
	}
	item, err := newItem(option)
	if err != nil {
		return nil, err
	}
	err = i.menu.Actions().Add(item.action)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (i windowsItem) AddItems(options ...ItemOption) ([]Item, error) {
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

func newItem(option ItemOption) (windowsItem, error) {
	if option.Title == "" && option.Icon == "" {
		return windowsItem{action: walk.NewSeparatorAction()}, nil
	}

	var (
		res windowsItem
		err error
	)
	if option.Children != nil {
		res.menu, err = walk.NewMenu()
		if err != nil {
			return res, err
		}
		res.action = walk.NewMenuAction(res.menu)

		_, err = res.AddItems(option.Children...)
		if err != nil {
			return res, err
		}
	} else {
		res.action = walk.NewAction()

	}
	if option.Title != "" {
		if err = res.SetTitle(option.Title); err != nil {
			return res, err
		}
	}

	if option.Disabled {
		if err = res.SetEnabled(false); err != nil {
			return res, err
		}
	}

	if option.Invisible {
		if err = res.SetVisible(false); err != nil {
			return res, err
		}
	}

	if option.Checkable {
		err = res.SetCheckable(option.Checkable)
		if err != nil {
			return res, err
		}
		if option.Checked {
			if err = res.SetChecked(option.Checked); err != nil {
				return res, err
			}
		}
	}

	if option.Icon != "" {
		if err := res.SetIcon(option.Icon); err != nil {
			return res, err
		}
	}

	if option.Callback != nil {
		if err := res.SetCallback(option.Callback); err != nil {
			return res, err
		}
	}

	return res, nil
}
