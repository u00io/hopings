package centerpanel

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/u00io/hopings/system"
	"github.com/u00io/nuiforms/ui"
)

type CenterPanel struct {
	ui.Widget

	mtx    sync.Mutex
	result system.Result

	orderColumnIndex int
	orderAsc         bool

	tableResults *ui.Table
}

func NewCenterPanel() *CenterPanel {
	var c CenterPanel
	c.InitWidget()
	c.tableResults = ui.NewTable()
	curstomWidgets := map[string]ui.Widgeter{
		"tableresults": c.tableResults,
	}
	c.SetLayout(`
		<column>
			<label id="lblIP" text="---" />
			<label id="lblResults" text="---"/>
			<widget id="tableresults" />
		</column>
	`, &c, curstomWidgets)

	c.orderColumnIndex = 1

	c.tableResults.SetColumnCount(4)
	c.tableResults.SetColumnName(0, "Index")
	c.tableResults.SetColumnName(1, "IP Address")
	c.tableResults.SetColumnName(2, "Time, ms")
	c.tableResults.SetColumnName(3, "Country")

	c.tableResults.SetColumnWidth(0, 70)
	c.tableResults.SetColumnWidth(1, 200)
	c.tableResults.SetColumnWidth(2, 120)
	c.tableResults.SetColumnWidth(3, 180)

	go c.thUpdateData()

	c.AddTimer(100, c.OnTimerUpdate)

	return &c
}

func (c *CenterPanel) thUpdateData() {
	for {
		result := system.Instance.GetResult()
		c.mtx.Lock()
		c.result = result
		c.mtx.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func (c *CenterPanel) HandleSystemEvent(event system.Event) {
	if event.Name == "update" {
		c.updateData()
	}
}

func (c *CenterPanel) OnTimerUpdate() {
	c.updateData()
}

func (c *CenterPanel) updateData() {
	r := c.result

	lblIP, ok := c.FindWidgetByName("lblIP").(*ui.Label)
	if ok {
		lblIP.SetText(fmt.Sprintf("Target IP: %s | Country: %s", r.IP, r.Country))
	}

	lblResults, ok := c.FindWidgetByName("lblResults").(*ui.Label)
	if ok {
		lblResults.SetText(fmt.Sprintf("Status: %s | Hops: %d", r.Status, len(r.Hops)))
	}

	c.tableResults.SetRowCount(len(r.Hops))
	for i, hop := range r.Hops {
		// INDEX
		c.tableResults.SetCellText2(i, 0, fmt.Sprint(i))

		// IP ADDRESS
		c.tableResults.SetCellText2(i, 1, hop.IP)

		// TIME MS
		c.tableResults.SetCellText2(i, 2, fmt.Sprint(hop.TimeMs))
		c.tableResults.SetCellColor(i, 2, color.RGBA{100, 100, 100, 255})

		// COUNTRY
		country, err := system.GetCountryByIP(hop.IP)
		if err != nil {
			country = ""
		}
		c.tableResults.SetCellText2(i, 3, country)
	}
}
