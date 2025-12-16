package bottompanel

import (
	"time"

	"github.com/u00io/hopings/system"
	"github.com/u00io/nuiforms/ui"
)

type BottomPanel struct {
	ui.Widget

	dtButtonCopyChanged time.Time
}

func NewBottomPanel() *BottomPanel {
	var c BottomPanel
	c.InitWidget()
	c.SetLayout(`
		<row>
			<button id="btnCopy" text="Copy" onclick="OnCopyClicked" />
			<hspacer />
			<button text="About" onclick="OnAboutClicked" />
		</row>
	`, &c, nil)

	c.SetElevation(5)

	c.AddTimer(1000, c.OnTimerUpdate)
	return &c
}

func (c *BottomPanel) HandleSystemEvent(event system.Event) {
}

func (c *BottomPanel) OnAboutClicked() {
	ui.ShowAboutDialog("About", "Hopings v0.2.2", "", "", "GeoLite2 data © MaxMind")
}

func (c *BottomPanel) OnCopyClicked() {
	system.Instance.CopyResultsToClipboard()
	btnCopy, ok := c.FindWidgetByName("btnCopy").(*ui.Button)
	if ok {
		btnCopy.SetText("Copied!")
		btnCopy.SetBackgroundColor(ui.ColorFromHex("#589c15ff"))
		c.dtButtonCopyChanged = time.Now()
	}
}

func (c *BottomPanel) OnTimerUpdate() {
	if time.Since(c.dtButtonCopyChanged) > time.Second {
		btnCopy, ok := c.FindWidgetByName("btnCopy").(*ui.Button)
		if ok {
			btnCopy.SetText("Copy")
			btnCopy.SetBackgroundColor(nil)
		}
	}
}
