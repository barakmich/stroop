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
	if err := g.SetKeybinding("", 'o', gocui.ModNone, c.openComment); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, c.openComment); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, c.clearOpen); err != nil {
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
	var err error
	switch c.mode {
	case "comment":
		err = c.curComment.MoveUp(g)
	default:
		err = c.columns[c.columnSelected].MoveUp(g)
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) moveDown(g *gocui.Gui, _ *gocui.View) error {
	var err error
	switch c.mode {
	case "comment":
		err = c.curComment.MoveDown(g)
	default:
		err = c.columns[c.columnSelected].MoveDown(g)
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) moveLeft(g *gocui.Gui, _ *gocui.View) error {
	var err error
	switch c.mode {
	case "comment":
	default:
		if c.columnSelected == 0 {
			return nil
		}
		err = c.columns[c.columnSelected].Deactivate(g)
		if err != nil {
			return err
		}
		c.columnSelected--
		err = c.columns[c.columnSelected].Activate(g)
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) moveRight(g *gocui.Gui, _ *gocui.View) error {
	var err error
	switch c.mode {
	case "comment":
	default:
		if c.columnSelected+1 == len(c.columns) {
			return nil
		}
		err := c.columns[c.columnSelected].Deactivate(g)
		if err != nil {
			return err
		}
		c.columnSelected++
		err = c.columns[c.columnSelected].Activate(g)
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *controller) openComment(g *gocui.Gui, _ *gocui.View) error {
	if c.mode == "comment" {
		return nil
	}
	issue, err := c.columns[c.columnSelected].CurrentIssue()
	if err != nil {
		return err
	}
	cv, err := CreateCommentView(g, c.client.GetCommentsForIssue(issue))
	if err != nil {
		return err
	}
	c.curComment = cv
	c.mode = "comment"
	return nil
}

func (c *controller) clearOpen(g *gocui.Gui, _ *gocui.View) error {
	if c.mode == "comment" {
		err := c.curComment.DeleteCommentView(g)
		if err != nil {
			return err
		}
		c.curComment = nil
		c.mode = ""
	}
	return nil
}
