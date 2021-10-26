package gelf

import (
	"encoding/json"
)

/*
{
   "_index":"logstash-2016.06.14",
   "_type":"docker",
   "_id":"AVVOP87XaNgHqqvkRf_V",
   "_score":null,
   "_source":{
      "version":"1.1",
      "host":"hamakabi",
      "level":6,
      "@version":"1",
      "@timestamp":"2016-06-14T09:30:51.854Z",
      "source_host":"172.17.0.1",
      "message":"2016-06-14 09:30:50 UTC LOG: incomplete startup packet",
      "command":"/bin/sh -c rm -f /var/log/postgresql/postgresql-9.3-main.log \t&& service postgresql start \t&& tail -f /var/log/postgresql/postgresql-9.3-main.log",
      "container_id":"2397fda5076279766a84950b0760844b5470259cae",
      "container_name":"dbgelf",
      "created":"2016-06-14T09:30:47.481540705Z",
      "image_id":"sha256:b4dc0e3d48524c6908170394e",
      "image_name":"pgsql:dev",
      "tag":"",
      "type":"docker"
   },
   "fields":{
      "created":[
         1465896647481
      ],
      "@timestamp":[
         1465896651854
      ]
   },
   "sort":[
      1465896651854
   ]
}
*/

type udpChunkedMsgSt struct {
	l      int
	sq     int
	chunks [][]byte
}

type messageSt struct {
	Version  string                 `json:"version"`
	Host     string                 `json:"host"`
	Short    string                 `json:"short_message"`
	Full     string                 `json:"full_message,omitempty"`
	TimeUnix float64                `json:"timestamp"`
	Level    int32                  `json:"level,omitempty"`
	Facility string                 `json:"facility,omitempty"`
	Extra    map[string]interface{} `json:"-"`
	RawExtra json.RawMessage        `json:"-"`
}
