package charts

type Bar struct {
	*Charts
}

func NewBar() *Bar {
	return &Bar{
		NewCharts(),
	}
}

func (this *Bar) Component() string {
	return "charts-bar"
}

// https://v-charts.js.org/#/bar

// chartSettings

// 指标维度
// metrics
func (this *Bar) Metrics(metrics []string) {
	this.WithSettings("metrics", metrics)
}

// dimension
func (this *Bar) Dimension(dimension []string) {
	this.WithSettings("dimension", dimension)
}

// 设置别名
// labelMap: {
//          'PV': '访问用户',
//          'Order': '下单用户'
//        },
func (this *Bar) LabelMap(maps map[string]interface{}) {
	this.WithSettings("labelMap", maps)
}

// https://v-charts.js.org/#/bar?id=%e8%ae%be%e7%bd%aelegend%e5%88%ab%e5%90%8d
func (this *Bar) LegendName(maps map[string]interface{}) {
	this.WithSettings("legendName", maps)
}

// https://v-charts.js.org/#/bar?id=%e5%a0%86%e5%8f%a0%e6%9d%a1%e5%bd%a2%e5%9b%be
func (this *Bar) Stack(columns []string) {
	stack := map[string]interface{}{}
	stack["stack"] = columns
	this.WithSettings("stack", stack)
}

// https://v-charts.js.org/#/bar?id=%e8%ae%be%e7%bd%ae%e7%ba%b5%e8%bd%b4%e4%b8%ba%e8%bf%9e%e7%bb%ad%e7%9a%84%e6%95%b0%e5%80%bc%e8%bd%b4
func (this *Bar) YAxisType() {
	this.WithSettings("yAxisType", "value")
}
