package toppanel

import (
	"github.com/u00io/hopings/system"
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nuiforms/ui"
)

type TopPanel struct {
	ui.Widget
}

func NewTopPanel() *TopPanel {
	var c TopPanel
	c.InitWidget()
	c.SetElevation(5)
	c.SetLayout(`
	<column>
		<row>
			<textbox id="txtTarget" text="example.com"/>
			<button id="btnStart" text="Start" onclick="OnStartClick"/>
			<button id="btnStop" text="Stop" onclick="OnStopClick"/>
		</row>
	</column>
	`, &c, nil)

	txtTarget, ok := c.FindWidgetByName("txtTarget").(*ui.TextBox)
	if ok {
		txtTarget.SetOnTextBoxKeyDown(func() {
			ev := ui.CurrentEvent().Parameter.(*ui.EventTextboxKeyDown)
			if ev.Key == nuikey.KeyEnter {
				c.OnStartClick()
				ev.Processed = true
				return
			}
			if ev.Key == nuikey.KeyEsc {
				c.OnStopClick()
				ev.Processed = true
				return
			}
		})
	}

	c.AddTimer(100, c.timerUpdate)
	return &c
}

func (c *TopPanel) timerUpdate() {
	c.updateButtons()
}

func (c *TopPanel) updateButtons() {
	btnStart, ok := c.FindWidgetByName("btnStart").(*ui.Button)
	if !ok {
		return
	}
	if system.Instance.IsStarted() {
		btnStart.SetRole("")
		btnStart.SetEnabled(false)
	} else {
		btnStart.SetRole("primary")
		btnStart.SetEnabled(true)
	}

	btnStop, ok := c.FindWidgetByName("btnStop").(*ui.Button)
	if !ok {
		return
	}
	if system.Instance.IsStarted() {
		btnStop.SetRole("primary")
		btnStop.SetEnabled(true)
	} else {
		btnStop.SetRole("")
		btnStop.SetEnabled(false)
	}
}

func (c *TopPanel) OnStartClick() {
	txtTarget, ok := c.FindWidgetByName("txtTarget").(*ui.TextBox)
	if !ok {
		return
	}
	system.Instance.Run(txtTarget.Text())
	c.updateButtons()
}

func (c *TopPanel) OnStopClick() {
	system.Instance.Abort()
	c.updateButtons()
}

func (c *TopPanel) HandleSystemEvent(event system.Event) {
}
