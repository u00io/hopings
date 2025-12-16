package bottompanel

import (
	"github.com/u00io/hopings/system"
	"github.com/u00io/nuiforms/ui"
)

type BottomPanel struct {
	ui.Widget
}

func NewBottomPanel() *BottomPanel {
	var c BottomPanel
	c.InitWidget()
	c.SetLayout(`
		<row>
			<hspacer />
			<button text="About" onclick="OnAboutClicked" />
		</row>
	`, &c, nil)

	c.SetElevation(5)
	return &c
}

func (c *BottomPanel) HandleSystemEvent(event system.Event) {
}

func (c *BottomPanel) OnAboutClicked() {
	ui.ShowAboutDialog("About", "Hopings v0.2.2", "", "", "GeoLite2 data © MaxMind")
}
