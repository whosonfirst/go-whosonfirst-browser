CWD=$(shell pwd)

cli:
	go build -o bin/server cmd/server/main.go

debug-tangram:
	go run -mod vendor cmd/server/main.go -map-provider tangram -nextzen-apikey $(APIKEY) -leaflet-enable-draw -javascript-at-eof

debug-tilepack:
	go run -mod vendor cmd/server/main.go -map-provider tangram -tilezen-enable-tilepack -tilezen-tilepack-path /usr/local/data/sf.db -leaflet-enable-draw

debug-protomaps:
	go run -mod vendor cmd/server/main.go -map-provider protomaps -protomaps-serve-tiles -protomaps-bucket-uri file://$(CWD)/fixtures -protomaps-database sfo -protomaps-paint-rules-uri file://$(CWD)/fixtures/protomaps.rules.paint.js -protomaps-label-rules-uri file://$(CWD)/fixtures/protomaps.rules.label.js -leaflet-enable-draw -javascript-at-eof

debug-leaflet:
	go run -mod vendor cmd/server/main.go -map-provider leaflet -leaflet-tile-url https://tile.openstreetmap.org/{z}/{x}/{y}.png -leaflet-enable-draw -javascript-at-eof
