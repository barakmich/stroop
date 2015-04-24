package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/coreos/pkg/capnslog"
	"github.com/jroimartin/gocui"
)

var log = capnslog.NewPackageLogger("github.com/barakmich/stroop", "main")

var debug = ""

type controller struct {
	client         *stroopClient
	columns        []*column
	defs           []*columnDef
	columnSelected int
	debugOn        bool
	statusleft     string
	statusright    string
}

func NewController(client *stroopClient) *controller {
	config := ReadConfig("$HOME/.stroop.conf")

	defs := config.Default
	if def, ok := config.Repos[client.remote.GithubName()]; ok {
		defs = def
	}
	c := &controller{
		client: client,
		defs:   defs,
	}
	return c
}

func main() {
	flag.Parse()
	capnslog.SetFormatter(capnslog.NewGlogFormatter(os.Stderr))

	client := startClient()
	client.MustDetectRemote()

	c := NewController(client)
	c.RefreshColumns()
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Fatal(err)
	}

	g.SetLayout(c.layout)

	if err := initKeybindings(g, c); err != nil {
		log.Fatal(err)
	}

	err := g.MainLoop()
	if err != nil && err != gocui.Quit {
		g.Close()
		log.Fatal(err)
	}
	g.Close()

}

func (c *controller) layout(g *gocui.Gui) error {
	g.BgColor = gocui.ColorDefault
	maxX, maxY := g.Size()
	if len(c.columns) == 0 {
		return errors.New("No columns defined")
	}
	colWidth := maxX / len(c.columns)
	for i, col := range c.columns {
		v, err := g.SetView(col.name, i*colWidth, 0, ((i+1)*colWidth)-1, maxY-2)
		if err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}
		}
		v.Clear()
		v.Write([]byte(col.name + ": \n\n"))
		err = col.CreateViews(g, v)
		if err != nil {
			return err
		}
	}
	status, err := g.SetView("statusbar", 0, maxY-2, maxX-1, maxY+1)
	if err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}

		status.Frame = false
		status.BgColor = gocui.ColorBlack
		status.FgColor = gocui.ColorDefault | gocui.AttrBold
	}
	status.Clear()
	if c.debugOn {
		status.Write([]byte(debug))
	} else {
		status.Write(c.statusBar(maxX))
	}
	return nil
}

func (c *controller) statusBar(width int) []byte {
	n := width - len(c.statusleft) - len(c.statusright) - 2
	status := fmt.Sprint(c.statusleft, strings.Repeat(" ", n), c.statusright)
	return []byte(status)
}

func (c *controller) RefreshColumns() {
	c.statusright = c.client.remote.GithubName()
	issues := c.client.GetIssues()
	columns := []*column{
		&column{
			name:     "Backlog",
			issues:   nil,
			isActive: true,
		},
	}
	backlog := columns[0]
	tagmap := make(map[string]*column)
	for _, def := range c.defs {
		col := &column{
			name:   def.Name,
			issues: nil,
		}
		columns = append(columns, col)
		tagmap[def.Tag] = col
	}

	for _, issue := range issues {
		foundTag := false
		for _, tag := range issue.Labels {
			if col, ok := tagmap[*tag.Name]; ok {
				col.issues = append(col.issues, issue)
				foundTag = true
				break
			}
		}
		if !foundTag {
			backlog.issues = append(backlog.issues, issue)
		}
	}
	c.columns = columns
}
