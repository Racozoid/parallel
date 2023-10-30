package charts

import (
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func dataToLineData(data []float64) []opts.LineData {
	var lineData []opts.LineData
	for _, value := range data {
		lineData = append(lineData, opts.LineData{Value: value})
	}

	return lineData
}

func DrawChartWith2Lines(data1, data2 []float64, XAxis []int, filename, title, XTitle, YTitle1, YTitle2 string) {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeInfographic,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: XTitle,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: title,
		}),
	)

	line.SetXAxis(XAxis).
		AddSeries(YTitle1, dataToLineData(data1)).
		AddSeries(YTitle2, dataToLineData(data2)).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	f, _ := os.Create(filename)
	_ = line.Render(f)
}
