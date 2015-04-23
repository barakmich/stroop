package main

import (
	"github.com/jroimartin/gocui"
)

func initKeybindings(g *gocui.Gui, c *controller) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'j', gocui.ModNone, c.moveDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'k', gocui.ModNone, c.moveUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'l', gocui.ModNone, c.moveRight); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'h', gocui.ModNone, c.moveLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error { c.debugOn = !c.debugOn; return nil }); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
}

func (c *controller) moveUp(g *gocui.Gui, _ *gocui.View) error {
	err := c.columns[c.columnSelected].MoveUp(g)
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) moveDown(g *gocui.Gui, _ *gocui.View) error {
	err := c.columns[c.columnSelected].MoveDown(g)
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) moveLeft(g *gocui.Gui, _ *gocui.View) error {
	if c.columnSelected == 0 {
		return nil
	}
	err := c.columns[c.columnSelected].Deactivate(g)
	if err != nil {
		return err
	}
	c.columnSelected--
	err = c.columns[c.columnSelected].Activate(g)
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) moveRight(g *gocui.Gui, _ *gocui.View) error {
	if c.columnSelected+1 == len(c.columns) {
		return nil
	}
	err := c.columns[c.columnSelected].Deactivate(g)
	if err != nil {
		return err
	}
	c.columnSelected++
	err = c.columns[c.columnSelected].Activate(g)
	if err != nil {
		return err
	}
	return nil
}
