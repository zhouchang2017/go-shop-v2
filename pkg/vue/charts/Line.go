package charts

// 折线图
type Line struct {
	*Charts
}

func NewLine() *Line {
	return &Line{
		NewCharts(),
	}
}

func (this *Line) Component() string {
	return "charts-line"
}

// https://v-charts.js.org/#/line
// chartSettings

// 指标维度
// metrics
func (this *Line) Metrics(metrics []string) {
	this.WithSettings("metrics", metrics)
}

// dimension
func (this *Line) Dimension(dimension []string) {
	this.WithSettings("dimension", dimension)
}

// 堆叠面积图
// area: true

// 设置别名
// labelMap: {
//          'PV': '访问用户',
//          'Order': '下单用户'
//        },
func (this *Line) LabelMap(maps map[string]interface{}) {
	this.WithSettings("labelMap", maps)
}

func (this *Line) LegendName(maps map[string]interface{}) {
	this.WithSettings("legendName", maps)
}

// 设置横轴为连续的数值轴
// xAxisType: 'value'
func (this *Line) XAxisTypeValue() {
	this.WithSettings("xAxisType", "value")
}

// 设置横轴为连续的时间轴
// xAxisType: 'time'
func (this *Line) XAxisTypeTime() {
	this.WithSettings("xAxisType", "time")
}

// extend

// 横坐标的倾斜
// 'xAxis.0.axisLabel.rotate': 45
func (this *Line) XLabelRotate(rotate int64) {
	this.WithExtend("xAxis.0.axisLabel.rotate", rotate)
}

// 显示指标数值
// series: {
//          label: {
//            normal: {
//              show: true
//            }
//          }
//        }
func (this *Line) LabelNormalShow() {
	this.WithExtend("series.label.normal.show", true)
}
