package centerpanel

import (
	_ "embed"
	"fmt"
	"image"
	"image/png"
	"strings"
	"sync"
	"time"

	"github.com/u00io/hopings/flags"
	"github.com/u00io/hopings/geomap"
	"github.com/u00io/hopings/system"
	"github.com/u00io/nuiforms/ui"
)

type CenterPanel struct {
	ui.Widget

	mtx    sync.Mutex
	result system.Result

	orderColumnIndex int
	orderAsc         bool

	strMapDigest string

	tableResults *ui.Table
}

//go:embed bluerect.png
var bluerect []byte

//go:embed redrect.png
var redrect []byte

//go:embed greenrect.png
var greenrect []byte

//go:embed grayrect.png
var grayrect []byte

func blueRect() image.Image {
	im, _ := png.Decode(strings.NewReader(string(bluerect)))
	return im
}

func redRect() image.Image {
	im, _ := png.Decode(strings.NewReader(string(redrect)))
	return im
}

func greenRect() image.Image {
	im, _ := png.Decode(strings.NewReader(string(greenrect)))
	return im
}

func grayRect() image.Image {
	im, _ := png.Decode(strings.NewReader(string(grayrect)))
	return im
}

func NewCenterPanel() *CenterPanel {
	var c CenterPanel
	c.InitWidget()
	c.tableResults = ui.NewTable()
	curstomWidgets := map[string]ui.Widgeter{
		"tableresults": c.tableResults,
	}
	c.SetLayout(`
	<row>
		<column>
			<row padding="0" spacing="0">
				<imagebox id="imgIP" />
				<space width="10" />
				<label id="lblIP" text="---" />
			</row>
			<row padding="0" spacing="0">
				<imagebox id="imgFlag" />
				<space width="10" />
				<label id="lblCountry" text="---" />
			</row>
			<row padding="0" spacing="0">
				<imagebox id="imgStatus" />
				<space width="10" />
				<label id="lblResults" text="---"/>
			</row>
			<widget id="tableresults" />
		</column>
		<imagebox id="imgMap" />
	</row>
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

	imgIP, ok := c.FindWidgetByName("imgIP").(*ui.ImageBox)
	if ok {
		imgIP.SetMinSize(32, 24)
		imgIP.SetMaxSize(32, 24)
		imgIP.SetSize(32, 24)
		imgIP.SetScaling(ui.ImageBoxScaleAdjustImageKeepAspectRatio)
		imgIP.SetImage(grayRect())
	}

	imgFlag, ok := c.FindWidgetByName("imgFlag").(*ui.ImageBox)
	if ok {
		imgFlag.SetMinSize(32, 24)
		imgFlag.SetMaxSize(32, 24)
		imgFlag.SetSize(32, 24)
		imgFlag.SetScaling(ui.ImageBoxScaleAdjustImageKeepAspectRatio)
		imgFlag.SetImage(grayRect())
	}

	imgStatus, ok := c.FindWidgetByName("imgStatus").(*ui.ImageBox)
	if ok {
		imgStatus.SetMinSize(32, 24)
		imgStatus.SetMaxSize(32, 24)
		imgStatus.SetSize(32, 24)
		imgStatus.SetScaling(ui.ImageBoxScaleAdjustImageKeepAspectRatio)
		imgStatus.SetImage(grayRect())
	}

	imgMap, ok := c.FindWidgetByName("imgMap").(*ui.ImageBox)
	if ok {
		imgMap.SetXExpandable(true)
		imgMap.SetYExpandable(true)
		//imgMap.SetMinSize(600, 400)
		//imgMap.SetMaxSize(32, 24)
		//imgMap.SetSize(32, 24)
		imgMap.SetScaling(ui.ImageBoxScaleAdjustImageKeepAspectRatio)
		imgMap.SetImage(grayRect())
	}

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

func (c *CenterPanel) updateMap(countries []string) {

	digest := strings.Join(countries, ",")
	if digest == c.strMapDigest {
		return
	}
	c.strMapDigest = digest

	settings := geomap.NewSettings()
	settings.Width = 1200
	settings.Height = 800
	settings.ShowAllCapitals = true
	settings.HighlightPath = []geomap.HighlightPoint{}

	for _, countryISO := range countries {
		if countryISO == "" {
			continue
		}
		settings.HighlightPath = append(settings.HighlightPath, geomap.HighlightPoint{
			CountryCode: countryISO,

			Style: geomap.DefaultMarkerStyle(),
		})
	}

	img, err := geomap.RenderMap(settings)
	if err != nil {
		panic(err)
	}

	imgMap, ok := c.FindWidgetByName("imgMap").(*ui.ImageBox)
	if ok {
		imgMap.SetImage(img)
	}

}

func (c *CenterPanel) OnTimerUpdate() {
	c.updateData()
}

func (c *CenterPanel) updateData() {
	r := c.result

	lblIP, ok := c.FindWidgetByName("lblIP").(*ui.Label)
	if ok {
		lblIP.SetText(fmt.Sprintf("Target IP: %s", r.IP))
	}

	lblCountry, ok := c.FindWidgetByName("lblCountry").(*ui.Label)
	if ok {
		lblCountry.SetText(fmt.Sprintf("Country: %s", r.CountryName))
	}

	imgFlag, ok := c.FindWidgetByName("imgFlag").(*ui.ImageBox)
	if ok {
		imgFlag.SetMinSize(32, 24)
		imgFlag.SetMaxSize(32, 24)
		imgFlag.SetSize(32, 24)
		imgFlag.SetScaling(ui.ImageBoxScaleAdjustImageKeepAspectRatio)
		if r.CountryISO == "" {
			imgFlag.SetImage(grayRect())
		} else {
			im, _ := flags.GetFlagImage(r.CountryISO)
			imgFlag.SetImage(im)
		}
	}

	lblResults, ok := c.FindWidgetByName("lblResults").(*ui.Label)
	if ok {
		lblResults.SetText(fmt.Sprintf("Status: %s | Hops: %d", r.Status, len(r.Hops)))
		if r.Status == "Finished" {
			lblResults.SetForegroundColor(ui.ColorFromHex("#23b423"))
			imgStatus, ok := c.FindWidgetByName("imgStatus").(*ui.ImageBox)
			if ok {
				imgStatus.SetImage(greenRect())
			}
		}
		if r.Status == "Running" {
			lblResults.SetForegroundColor(ui.ColorFromHex("#3bafe6"))
			imgStatus, ok := c.FindWidgetByName("imgStatus").(*ui.ImageBox)
			if ok {
				imgStatus.SetImage(blueRect())
			}
		}
		if strings.Contains(r.Status, "ERR:") {
			lblResults.SetForegroundColor(ui.ColorFromHex("#FF7777"))
			imgStatus, ok := c.FindWidgetByName("imgStatus").(*ui.ImageBox)
			if ok {
				imgStatus.SetImage(redRect())
			}
		}
	}

	countries := make([]string, 0)

	c.tableResults.SetRowCount(len(r.Hops))
	for i, hop := range r.Hops {
		// INDEX
		c.tableResults.SetCellText2(i, 0, fmt.Sprint(i))

		// IP ADDRESS
		c.tableResults.SetCellText2(i, 1, hop.IP)
		if hop.IP == "" {
			c.tableResults.SetCellText2(i, 1, "-- no reply -- ")
			c.tableResults.SetCellColor(i, 1, ui.ColorFromHex("#777777"))
		} else {
			col := ui.ColorFromHex("#FFFFFF")
			if hop.IP == r.IP {
				col = ui.ColorFromHex("#76e676")
			}
			c.tableResults.SetCellColor(i, 1, col)
		}

		// TIME MS
		if hop.TimeMs < 0 {
			c.tableResults.SetCellText2(i, 2, "")
			c.tableResults.SetCellColor(i, 2, ui.ColorFromHex("#777777"))
		} else {
			col := ui.ColorFromHex("#FFFFFF")

			// Highlight fast times
			if hop.TimeMs >= 0 && hop.TimeMs < 100 {
				col = ui.ColorFromHex("#76e676")
			} else if hop.TimeMs >= 100 && hop.TimeMs < 200 {
				col = ui.ColorFromHex("#f0f05f")
			} else if hop.TimeMs >= 200 {
				col = ui.ColorFromHex("#f3a375ff")
			}

			c.tableResults.SetCellText2(i, 2, fmt.Sprint(hop.TimeMs))
			c.tableResults.SetCellColor(i, 2, col)
		}

		// COUNTRY
		c.tableResults.SetCellText2(i, 3, hop.CountryName)

		if hop.CountryISO != "" {
			im, _ := flags.GetFlagImage(hop.CountryISO)
			c.tableResults.SetCellImage(i, 3, im, 24)
		} else {
			c.tableResults.SetCellImage(i, 3, nil, 0)
		}

		countries = append(countries, hop.CountryISO)

	}

	strTable := ""
	for colIndex := 0; colIndex < c.tableResults.ColumnCount(); colIndex++ {
		strTable += c.tableResults.ColumnName(colIndex) + "\t"
	}
	strTable += "\r\n"
	for rowIndex := 0; rowIndex < c.tableResults.RowCount(); rowIndex++ {
		strRow := ""
		for colIndex := 0; colIndex < c.tableResults.ColumnCount(); colIndex++ {
			strRow += c.tableResults.GetCellText2(rowIndex, colIndex) + "\t"
		}
		strTable += strRow + "\r\n"
	}

	system.Instance.SetResultTableText(strTable)

	c.updateMap(countries)
}
