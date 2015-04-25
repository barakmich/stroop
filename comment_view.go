package main

import (
	"fmt"
	"strconv"

	"github.com/jroimartin/gocui"
)

type commentView struct {
	ic *issueComment
}

func (c *commentView) viewName() string {
	return fmt.Sprintf("comment-%d", *c.ic.issue.Number)
}

func CreateCommentView(g *gocui.Gui, i *issueComment) (*commentView, error) {
	c := &commentView{
		ic: i,
	}
	maxX, maxY := g.Size()
	_, err := g.SetView(c.viewName(), 4, 4, maxX-4, maxY-4)
	if err != gocui.ErrorUnkView {
		return nil, err
	}
	debug = *i.issue.URL

	return c, c.updateView(g)
}

func (c *commentView) MoveDown(g *gocui.Gui) error {
	v, err := g.View(c.viewName())
	if err != nil {
		return err
	}
	x, y := v.Origin()
	v.SetOrigin(x, y+1)
	return nil
}

func (c *commentView) MoveUp(g *gocui.Gui) error {
	v, err := g.View(c.viewName())
	if err != nil {
		return err
	}
	x, y := v.Origin()
	if y == 0 {
		return nil
	}
	v.SetOrigin(x, y-1)
	return nil
}

func (c *commentView) DeleteCommentView(g *gocui.Gui) error {
	return g.DeleteView(fmt.Sprintf("comment-%d", *c.ic.issue.Number))
}

func (c *commentView) updateView(g *gocui.Gui) error {
	v, err := g.View(c.viewName())
	if err != nil {
		return err
	}
	w, _ := v.Size()
	v.Wrap = true
	v.Clear()
	v.Write([]byte(CenterLine(strconv.Itoa(*c.ic.issue.Number), w)))
	v.Write([]byte("\n"))
	v.Write([]byte(CenterLine(*c.ic.issue.Title, w)))
	v.Write([]byte("\n"))
	v.Write([]byte(fmt.Sprintf("Reporter: %s", *c.ic.issue.User.Login)))
	v.Write([]byte("\n"))
	v.Write([]byte(CleanBody(c.ic.issue.Body)))
	v.Write([]byte("\n\n"))
	for _, c := range c.ic.comments {
		v.Write([]byte("@" + *c.User.Login + ":\n"))
		v.Write([]byte(CleanBody(c.Body)))
		v.Write([]byte("\n\n"))
	}
	return nil
}
