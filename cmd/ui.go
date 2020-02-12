package cmd

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/spf13/cobra"
)

var (
	uiCmd = &cobra.Command{
		Use: "ui",
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := gocui.NewGui(gocui.OutputNormal, false)
			if err != nil {
				return err
			}
			defer g.Close()

			g.SetManagerFunc(layout)

			if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
				return err
			}

			if err := g.MainLoop(); err != nil && gocui.IsQuit(err) {
				return err
			}

			return nil
		},
	}
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
		if _, err := g.SetCurrentView("hello"); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
