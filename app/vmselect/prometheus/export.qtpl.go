// Code generated by qtc from "export.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line app/vmselect/prometheus/export.qtpl:1
package prometheus

//line app/vmselect/prometheus/export.qtpl:1
import (
	"bytes"
	"math"
	"strings"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/querytracer"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/storage"
	"github.com/valyala/quicktemplate"
)

//line app/vmselect/prometheus/export.qtpl:14
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line app/vmselect/prometheus/export.qtpl:14
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line app/vmselect/prometheus/export.qtpl:14
func StreamExportCSVLine(qw422016 *qt422016.Writer, xb *exportBlock, fieldNames []string) {
//line app/vmselect/prometheus/export.qtpl:15
	if len(xb.timestamps) == 0 || len(fieldNames) == 0 {
//line app/vmselect/prometheus/export.qtpl:15
		return
//line app/vmselect/prometheus/export.qtpl:15
	}
//line app/vmselect/prometheus/export.qtpl:16
	for i, timestamp := range xb.timestamps {
//line app/vmselect/prometheus/export.qtpl:17
		value := xb.values[i]

//line app/vmselect/prometheus/export.qtpl:18
		streamexportCSVField(qw422016, xb.mn, fieldNames[0], timestamp, value)
//line app/vmselect/prometheus/export.qtpl:19
		for _, fieldName := range fieldNames[1:] {
//line app/vmselect/prometheus/export.qtpl:19
			qw422016.N().S(`,`)
//line app/vmselect/prometheus/export.qtpl:21
			streamexportCSVField(qw422016, xb.mn, fieldName, timestamp, value)
//line app/vmselect/prometheus/export.qtpl:22
		}
//line app/vmselect/prometheus/export.qtpl:23
		qw422016.N().S(`
`)
//line app/vmselect/prometheus/export.qtpl:24
	}
//line app/vmselect/prometheus/export.qtpl:25
}

//line app/vmselect/prometheus/export.qtpl:25
func WriteExportCSVLine(qq422016 qtio422016.Writer, xb *exportBlock, fieldNames []string) {
//line app/vmselect/prometheus/export.qtpl:25
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:25
	StreamExportCSVLine(qw422016, xb, fieldNames)
//line app/vmselect/prometheus/export.qtpl:25
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:25
}

//line app/vmselect/prometheus/export.qtpl:25
func ExportCSVLine(xb *exportBlock, fieldNames []string) string {
//line app/vmselect/prometheus/export.qtpl:25
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:25
	WriteExportCSVLine(qb422016, xb, fieldNames)
//line app/vmselect/prometheus/export.qtpl:25
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:25
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:25
	return qs422016
//line app/vmselect/prometheus/export.qtpl:25
}

//line app/vmselect/prometheus/export.qtpl:27
func streamexportCSVField(qw422016 *qt422016.Writer, mn *storage.MetricName, fieldName string, timestamp int64, value float64) {
//line app/vmselect/prometheus/export.qtpl:28
	if fieldName == "__value__" {
//line app/vmselect/prometheus/export.qtpl:29
		qw422016.N().F(value)
//line app/vmselect/prometheus/export.qtpl:30
		return
//line app/vmselect/prometheus/export.qtpl:31
	}
//line app/vmselect/prometheus/export.qtpl:32
	if fieldName == "__timestamp__" {
//line app/vmselect/prometheus/export.qtpl:33
		qw422016.N().DL(timestamp)
//line app/vmselect/prometheus/export.qtpl:34
		return
//line app/vmselect/prometheus/export.qtpl:35
	}
//line app/vmselect/prometheus/export.qtpl:36
	if strings.HasPrefix(fieldName, "__timestamp__:") {
//line app/vmselect/prometheus/export.qtpl:37
		timeFormat := fieldName[len("__timestamp__:"):]

//line app/vmselect/prometheus/export.qtpl:38
		switch timeFormat {
//line app/vmselect/prometheus/export.qtpl:39
		case "unix_s":
//line app/vmselect/prometheus/export.qtpl:40
			qw422016.N().DL(timestamp / 1000)
//line app/vmselect/prometheus/export.qtpl:41
		case "unix_ms":
//line app/vmselect/prometheus/export.qtpl:42
			qw422016.N().DL(timestamp)
//line app/vmselect/prometheus/export.qtpl:43
		case "unix_ns":
//line app/vmselect/prometheus/export.qtpl:44
			qw422016.N().DL(timestamp * 1e6)
//line app/vmselect/prometheus/export.qtpl:45
		case "rfc3339":
//line app/vmselect/prometheus/export.qtpl:47
			bb := quicktemplate.AcquireByteBuffer()
			bb.B = time.Unix(timestamp/1000, (timestamp%1000)*1e6).AppendFormat(bb.B[:0], time.RFC3339)

//line app/vmselect/prometheus/export.qtpl:50
			qw422016.N().Z(bb.B)
//line app/vmselect/prometheus/export.qtpl:52
			quicktemplate.ReleaseByteBuffer(bb)

//line app/vmselect/prometheus/export.qtpl:54
		default:
//line app/vmselect/prometheus/export.qtpl:55
			if strings.HasPrefix(timeFormat, "custom:") {
//line app/vmselect/prometheus/export.qtpl:57
				layout := timeFormat[len("custom:"):]
				bb := quicktemplate.AcquireByteBuffer()
				bb.B = time.Unix(timestamp/1000, (timestamp%1000)*1e6).AppendFormat(bb.B[:0], layout)

//line app/vmselect/prometheus/export.qtpl:61
				if bytes.ContainsAny(bb.B, `"`+",\n") {
//line app/vmselect/prometheus/export.qtpl:62
					qw422016.E().QZ(bb.B)
//line app/vmselect/prometheus/export.qtpl:63
				} else {
//line app/vmselect/prometheus/export.qtpl:64
					qw422016.N().Z(bb.B)
//line app/vmselect/prometheus/export.qtpl:65
				}
//line app/vmselect/prometheus/export.qtpl:67
				quicktemplate.ReleaseByteBuffer(bb)

//line app/vmselect/prometheus/export.qtpl:69
			} else {
//line app/vmselect/prometheus/export.qtpl:69
				qw422016.N().S(`Unsupported timeFormat=`)
//line app/vmselect/prometheus/export.qtpl:70
				qw422016.N().S(timeFormat)
//line app/vmselect/prometheus/export.qtpl:71
			}
//line app/vmselect/prometheus/export.qtpl:72
		}
//line app/vmselect/prometheus/export.qtpl:73
		return
//line app/vmselect/prometheus/export.qtpl:74
	}
//line app/vmselect/prometheus/export.qtpl:75
	v := mn.GetTagValue(fieldName)

//line app/vmselect/prometheus/export.qtpl:76
	if bytes.ContainsAny(v, `"`+",\n") {
//line app/vmselect/prometheus/export.qtpl:77
		qw422016.N().QZ(v)
//line app/vmselect/prometheus/export.qtpl:78
	} else {
//line app/vmselect/prometheus/export.qtpl:79
		qw422016.N().Z(v)
//line app/vmselect/prometheus/export.qtpl:80
	}
//line app/vmselect/prometheus/export.qtpl:81
}

//line app/vmselect/prometheus/export.qtpl:81
func writeexportCSVField(qq422016 qtio422016.Writer, mn *storage.MetricName, fieldName string, timestamp int64, value float64) {
//line app/vmselect/prometheus/export.qtpl:81
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:81
	streamexportCSVField(qw422016, mn, fieldName, timestamp, value)
//line app/vmselect/prometheus/export.qtpl:81
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:81
}

//line app/vmselect/prometheus/export.qtpl:81
func exportCSVField(mn *storage.MetricName, fieldName string, timestamp int64, value float64) string {
//line app/vmselect/prometheus/export.qtpl:81
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:81
	writeexportCSVField(qb422016, mn, fieldName, timestamp, value)
//line app/vmselect/prometheus/export.qtpl:81
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:81
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:81
	return qs422016
//line app/vmselect/prometheus/export.qtpl:81
}

//line app/vmselect/prometheus/export.qtpl:83
func StreamExportPrometheusLine(qw422016 *qt422016.Writer, xb *exportBlock) {
//line app/vmselect/prometheus/export.qtpl:84
	if len(xb.timestamps) == 0 {
//line app/vmselect/prometheus/export.qtpl:84
		return
//line app/vmselect/prometheus/export.qtpl:84
	}
//line app/vmselect/prometheus/export.qtpl:85
	bb := quicktemplate.AcquireByteBuffer()

//line app/vmselect/prometheus/export.qtpl:86
	writeprometheusMetricName(bb, xb.mn)

//line app/vmselect/prometheus/export.qtpl:87
	for i, ts := range xb.timestamps {
//line app/vmselect/prometheus/export.qtpl:88
		qw422016.N().Z(bb.B)
//line app/vmselect/prometheus/export.qtpl:88
		qw422016.N().S(` `)
//line app/vmselect/prometheus/export.qtpl:89
		qw422016.N().F(xb.values[i])
//line app/vmselect/prometheus/export.qtpl:89
		qw422016.N().S(` `)
//line app/vmselect/prometheus/export.qtpl:90
		qw422016.N().DL(ts)
//line app/vmselect/prometheus/export.qtpl:90
		qw422016.N().S(`
`)
//line app/vmselect/prometheus/export.qtpl:91
	}
//line app/vmselect/prometheus/export.qtpl:92
	quicktemplate.ReleaseByteBuffer(bb)

//line app/vmselect/prometheus/export.qtpl:93
}

//line app/vmselect/prometheus/export.qtpl:93
func WriteExportPrometheusLine(qq422016 qtio422016.Writer, xb *exportBlock) {
//line app/vmselect/prometheus/export.qtpl:93
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:93
	StreamExportPrometheusLine(qw422016, xb)
//line app/vmselect/prometheus/export.qtpl:93
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:93
}

//line app/vmselect/prometheus/export.qtpl:93
func ExportPrometheusLine(xb *exportBlock) string {
//line app/vmselect/prometheus/export.qtpl:93
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:93
	WriteExportPrometheusLine(qb422016, xb)
//line app/vmselect/prometheus/export.qtpl:93
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:93
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:93
	return qs422016
//line app/vmselect/prometheus/export.qtpl:93
}

//line app/vmselect/prometheus/export.qtpl:95
func StreamExportJSONLine(qw422016 *qt422016.Writer, xb *exportBlock) {
//line app/vmselect/prometheus/export.qtpl:96
	if len(xb.timestamps) == 0 {
//line app/vmselect/prometheus/export.qtpl:96
		return
//line app/vmselect/prometheus/export.qtpl:96
	}
//line app/vmselect/prometheus/export.qtpl:96
	qw422016.N().S(`{"metric":`)
//line app/vmselect/prometheus/export.qtpl:98
	streammetricNameObject(qw422016, xb.mn)
//line app/vmselect/prometheus/export.qtpl:98
	qw422016.N().S(`,"values":[`)
//line app/vmselect/prometheus/export.qtpl:100
	if len(xb.values) > 0 {
//line app/vmselect/prometheus/export.qtpl:101
		values := xb.values

//line app/vmselect/prometheus/export.qtpl:102
		streamconvertValueToSpecialJSON(qw422016, values[0])
//line app/vmselect/prometheus/export.qtpl:103
		values = values[1:]

//line app/vmselect/prometheus/export.qtpl:104
		for _, v := range values {
//line app/vmselect/prometheus/export.qtpl:104
			qw422016.N().S(`,`)
//line app/vmselect/prometheus/export.qtpl:105
			streamconvertValueToSpecialJSON(qw422016, v)
//line app/vmselect/prometheus/export.qtpl:106
		}
//line app/vmselect/prometheus/export.qtpl:107
	}
//line app/vmselect/prometheus/export.qtpl:107
	qw422016.N().S(`],"timestamps":[`)
//line app/vmselect/prometheus/export.qtpl:110
	if len(xb.timestamps) > 0 {
//line app/vmselect/prometheus/export.qtpl:111
		timestamps := xb.timestamps

//line app/vmselect/prometheus/export.qtpl:112
		qw422016.N().DL(timestamps[0])
//line app/vmselect/prometheus/export.qtpl:113
		timestamps = timestamps[1:]

//line app/vmselect/prometheus/export.qtpl:114
		for _, ts := range timestamps {
//line app/vmselect/prometheus/export.qtpl:114
			qw422016.N().S(`,`)
//line app/vmselect/prometheus/export.qtpl:115
			qw422016.N().DL(ts)
//line app/vmselect/prometheus/export.qtpl:116
		}
//line app/vmselect/prometheus/export.qtpl:117
	}
//line app/vmselect/prometheus/export.qtpl:117
	qw422016.N().S(`]}`)
//line app/vmselect/prometheus/export.qtpl:119
	qw422016.N().S(`
`)
//line app/vmselect/prometheus/export.qtpl:120
}

//line app/vmselect/prometheus/export.qtpl:120
func WriteExportJSONLine(qq422016 qtio422016.Writer, xb *exportBlock) {
//line app/vmselect/prometheus/export.qtpl:120
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:120
	StreamExportJSONLine(qw422016, xb)
//line app/vmselect/prometheus/export.qtpl:120
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:120
}

//line app/vmselect/prometheus/export.qtpl:120
func ExportJSONLine(xb *exportBlock) string {
//line app/vmselect/prometheus/export.qtpl:120
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:120
	WriteExportJSONLine(qb422016, xb)
//line app/vmselect/prometheus/export.qtpl:120
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:120
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:120
	return qs422016
//line app/vmselect/prometheus/export.qtpl:120
}

//line app/vmselect/prometheus/export.qtpl:122
func StreamExportPromAPILine(qw422016 *qt422016.Writer, xb *exportBlock) {
//line app/vmselect/prometheus/export.qtpl:122
	qw422016.N().S(`{"metric":`)
//line app/vmselect/prometheus/export.qtpl:124
	streammetricNameObject(qw422016, xb.mn)
//line app/vmselect/prometheus/export.qtpl:124
	qw422016.N().S(`,"values":`)
//line app/vmselect/prometheus/export.qtpl:125
	streamvaluesWithTimestamps(qw422016, xb.values, xb.timestamps)
//line app/vmselect/prometheus/export.qtpl:125
	qw422016.N().S(`}`)
//line app/vmselect/prometheus/export.qtpl:127
}

//line app/vmselect/prometheus/export.qtpl:127
func WriteExportPromAPILine(qq422016 qtio422016.Writer, xb *exportBlock) {
//line app/vmselect/prometheus/export.qtpl:127
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:127
	StreamExportPromAPILine(qw422016, xb)
//line app/vmselect/prometheus/export.qtpl:127
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:127
}

//line app/vmselect/prometheus/export.qtpl:127
func ExportPromAPILine(xb *exportBlock) string {
//line app/vmselect/prometheus/export.qtpl:127
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:127
	WriteExportPromAPILine(qb422016, xb)
//line app/vmselect/prometheus/export.qtpl:127
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:127
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:127
	return qs422016
//line app/vmselect/prometheus/export.qtpl:127
}

//line app/vmselect/prometheus/export.qtpl:129
func StreamExportPromAPIHeader(qw422016 *qt422016.Writer) {
//line app/vmselect/prometheus/export.qtpl:129
	qw422016.N().S(`{"status":"success","data":{"resultType":"matrix","result":[`)
//line app/vmselect/prometheus/export.qtpl:135
}

//line app/vmselect/prometheus/export.qtpl:135
func WriteExportPromAPIHeader(qq422016 qtio422016.Writer) {
//line app/vmselect/prometheus/export.qtpl:135
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:135
	StreamExportPromAPIHeader(qw422016)
//line app/vmselect/prometheus/export.qtpl:135
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:135
}

//line app/vmselect/prometheus/export.qtpl:135
func ExportPromAPIHeader() string {
//line app/vmselect/prometheus/export.qtpl:135
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:135
	WriteExportPromAPIHeader(qb422016)
//line app/vmselect/prometheus/export.qtpl:135
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:135
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:135
	return qs422016
//line app/vmselect/prometheus/export.qtpl:135
}

//line app/vmselect/prometheus/export.qtpl:137
func StreamExportPromAPIFooter(qw422016 *qt422016.Writer, qt *querytracer.Tracer) {
//line app/vmselect/prometheus/export.qtpl:137
	qw422016.N().S(`]}`)
//line app/vmselect/prometheus/export.qtpl:141
	qt.Donef("export format=promapi")

//line app/vmselect/prometheus/export.qtpl:143
	streamdumpQueryTrace(qw422016, qt)
//line app/vmselect/prometheus/export.qtpl:143
	qw422016.N().S(`}`)
//line app/vmselect/prometheus/export.qtpl:145
}

//line app/vmselect/prometheus/export.qtpl:145
func WriteExportPromAPIFooter(qq422016 qtio422016.Writer, qt *querytracer.Tracer) {
//line app/vmselect/prometheus/export.qtpl:145
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:145
	StreamExportPromAPIFooter(qw422016, qt)
//line app/vmselect/prometheus/export.qtpl:145
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:145
}

//line app/vmselect/prometheus/export.qtpl:145
func ExportPromAPIFooter(qt *querytracer.Tracer) string {
//line app/vmselect/prometheus/export.qtpl:145
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:145
	WriteExportPromAPIFooter(qb422016, qt)
//line app/vmselect/prometheus/export.qtpl:145
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:145
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:145
	return qs422016
//line app/vmselect/prometheus/export.qtpl:145
}

//line app/vmselect/prometheus/export.qtpl:147
func streamprometheusMetricName(qw422016 *qt422016.Writer, mn *storage.MetricName) {
//line app/vmselect/prometheus/export.qtpl:148
	qw422016.N().Z(mn.MetricGroup)
//line app/vmselect/prometheus/export.qtpl:149
	if len(mn.Tags) > 0 {
//line app/vmselect/prometheus/export.qtpl:149
		qw422016.N().S(`{`)
//line app/vmselect/prometheus/export.qtpl:151
		tags := mn.Tags

//line app/vmselect/prometheus/export.qtpl:152
		qw422016.N().Z(tags[0].Key)
//line app/vmselect/prometheus/export.qtpl:152
		qw422016.N().S(`=`)
//line app/vmselect/prometheus/export.qtpl:152
		qw422016.N().QZ(tags[0].Value)
//line app/vmselect/prometheus/export.qtpl:153
		tags = tags[1:]

//line app/vmselect/prometheus/export.qtpl:154
		for i := range tags {
//line app/vmselect/prometheus/export.qtpl:155
			tag := &tags[i]

//line app/vmselect/prometheus/export.qtpl:155
			qw422016.N().S(`,`)
//line app/vmselect/prometheus/export.qtpl:156
			qw422016.N().Z(tag.Key)
//line app/vmselect/prometheus/export.qtpl:156
			qw422016.N().S(`=`)
//line app/vmselect/prometheus/export.qtpl:156
			qw422016.N().QZ(tag.Value)
//line app/vmselect/prometheus/export.qtpl:157
		}
//line app/vmselect/prometheus/export.qtpl:157
		qw422016.N().S(`}`)
//line app/vmselect/prometheus/export.qtpl:159
	}
//line app/vmselect/prometheus/export.qtpl:160
}

//line app/vmselect/prometheus/export.qtpl:160
func writeprometheusMetricName(qq422016 qtio422016.Writer, mn *storage.MetricName) {
//line app/vmselect/prometheus/export.qtpl:160
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:160
	streamprometheusMetricName(qw422016, mn)
//line app/vmselect/prometheus/export.qtpl:160
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:160
}

//line app/vmselect/prometheus/export.qtpl:160
func prometheusMetricName(mn *storage.MetricName) string {
//line app/vmselect/prometheus/export.qtpl:160
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:160
	writeprometheusMetricName(qb422016, mn)
//line app/vmselect/prometheus/export.qtpl:160
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:160
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:160
	return qs422016
//line app/vmselect/prometheus/export.qtpl:160
}

//line app/vmselect/prometheus/export.qtpl:187
func streamconvertValueToSpecialJSON(qw422016 *qt422016.Writer, v float64) {
//line app/vmselect/prometheus/export.qtpl:188
	if math.IsNaN(v) {
//line app/vmselect/prometheus/export.qtpl:188
		qw422016.N().S(`"NaN"`)
//line app/vmselect/prometheus/export.qtpl:190
	} else if math.IsInf(v, 1) {
//line app/vmselect/prometheus/export.qtpl:190
		qw422016.N().S(`"Infinity"`)
//line app/vmselect/prometheus/export.qtpl:192
	} else if math.IsInf(v, -1) {
//line app/vmselect/prometheus/export.qtpl:192
		qw422016.N().S(`"-Infinity"`)
//line app/vmselect/prometheus/export.qtpl:194
	} else {
//line app/vmselect/prometheus/export.qtpl:195
		qw422016.N().F(v)
//line app/vmselect/prometheus/export.qtpl:196
	}
//line app/vmselect/prometheus/export.qtpl:197
}

//line app/vmselect/prometheus/export.qtpl:197
func writeconvertValueToSpecialJSON(qq422016 qtio422016.Writer, v float64) {
//line app/vmselect/prometheus/export.qtpl:197
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:197
	streamconvertValueToSpecialJSON(qw422016, v)
//line app/vmselect/prometheus/export.qtpl:197
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:197
}

//line app/vmselect/prometheus/export.qtpl:197
func convertValueToSpecialJSON(v float64) string {
//line app/vmselect/prometheus/export.qtpl:197
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:197
	writeconvertValueToSpecialJSON(qb422016, v)
//line app/vmselect/prometheus/export.qtpl:197
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:197
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:197
	return qs422016
//line app/vmselect/prometheus/export.qtpl:197
}
