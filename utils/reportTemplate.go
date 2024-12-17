package utils

func reportTemplate() string {
	return `
<strong>Daily Report - {{.ReportTime}}</strong>

<b>📈 Alert! [Crossed Threshold]</b>
<pre>
{{range .PriceGainers}}
{{.Ticker | printf "%-10s"}} Current Price: {{printf "%.2f" .CurrentPrice.Float64}} BP: {{printf "%.2f" .BuyingPrice}}
{{end}}
</pre>

<b>📈 Potential Gainers (SMA50)</b>
<pre>
{{range .SMA50Gainers}}
{{.Ticker | printf "%-10s"}} Current Price: {{printf "%.2f" .CurrentPrice.Float64}}  SMA50: {{printf "%.2f" .Sma50.Float64}}
{{end}}
</pre>

<b>📉 Potential Losers (SMA50)</b>
<pre>
{{range .SMA50Losers}}
{{.Ticker | printf "%-10s"}} Current Price: {{printf "%.2f" .CurrentPrice.Float64}}  SMA50: {{printf "%.2f" .Sma50.Float64}}
{{end}}
</pre>
`
}
