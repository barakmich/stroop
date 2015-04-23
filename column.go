package main

import (
	"fmt"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/jroimartin/gocui"
)

const (
	LinesPerEntry = 5
)

type column struct {
	name      string
	issues    []github.Issue
	top       int
	maxIssues int
	// -1 is unselected
	selection int
	isActive  bool
}

type columnDef struct {
	Name string `json:"column"`
	Tag  string `json:"tag"`
}

func (c *column) CreateViews(g *gocui.Gui, colv *gocui.View) error {
	if c.isActive {
		colv.FgColor = gocui.ColorCyan | gocui.AttrBold
	} else {
		colv.FgColor = gocui.ColorDefault
	}
	x, y, maxX, maxY, err := g.ViewPosition(colv.Name())
	y = y + 2
	if err != nil {
		return err
	}
	maxIssues := maxY / LinesPerEntry
	c.maxIssues = maxIssues
	for i := 0; i < maxIssues; i++ {
		v, err := g.SetView(fmt.Sprintf("col-%s-%d", c.name, i),
			x, y+(i*LinesPerEntry), maxX, y+((i+1)*LinesPerEntry))
		if err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}
		}
		v.SelBgColor = gocui.ColorRed
		v.Frame = false
		v.Wrap = true
	}
	return c.redraw(g)
}

func (c *column) viewName(i int) string {
	return fmt.Sprintf("col-%s-%d", c.name, i)
}

func (c *column) MoveUp(g *gocui.Gui) error {
	if c.selection == -1 {
		return nil
	}
	if c.top == 0 && c.selection == 0 {
		return nil
	}
	if c.top == c.selection {
		c.top--
	}
	c.selection--
	return c.redraw(g)
}

func (c *column) Activate(g *gocui.Gui) error {
	c.isActive = true
	return c.redraw(g)
}

func (c *column) Deactivate(g *gocui.Gui) error {
	c.isActive = false
	return c.redraw(g)
}

func (c *column) issueDetails(issue *github.Issue) string {
	if (issue.Comments == nil || *issue.Comments == 0) && issue.Assignee == nil {
		return ""
	}
	n := 0
	if issue.Comments != nil {
		n = *issue.Comments
	}
	s := strconv.Itoa(n)
	if issue.Assignee != nil {
		s += fmt.Sprint(":@", *issue.Assignee.Login)
	}
	return fmt.Sprintf("(%s)", s)
}

func (c *column) redraw(g *gocui.Gui) error {
	for i := 0; i < c.maxIssues && c.top+i < len(c.issues); i++ {
		v, err := g.View(c.viewName(i))
		if err != nil {
			return err
		}
		v.Clear()
		issue := c.issues[c.top+i]
		v.Write([]byte(strconv.Itoa(*issue.Number)))
		v.Write([]byte(fmt.Sprintf(": %s\n", c.issueDetails(&issue))))
		v.Write([]byte(*issue.Title))
		v.Write([]byte("\n"))
		if c.isActive && c.selection == c.top+i {
			v.BgColor = gocui.ColorRed
		} else {
			v.BgColor = gocui.ColorDefault
		}
	}
	return nil
}

func (c *column) MoveDown(g *gocui.Gui) error {
	if c.selection == -1 {
		return nil
	}
	if c.selection == len(c.issues)-1 {
		return nil
	}
	if c.top+c.maxIssues-1 == c.selection {
		c.top++
	}
	c.selection++
	debug = fmt.Sprintf("top: %d, select: %d, active: %v", c.top, c.selection, c.isActive)
	return c.redraw(g)
}
