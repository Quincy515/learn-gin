package resource

var (
  {{range $k,$v:=. }}
        {{$k}}=`{{Gzip $v}}`
   {{end}}
)