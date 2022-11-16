module github.com/whosonfirst/go-whosonfirst-browser/v5

go 1.19

// Note: until I can figure out why v1.6.0 doesn't work
// work with go-sfomuseum-pmtiles we are pegged at v1.5.1
// github.com/protomaps/go-pmtiles v1.6.0
// require github.com/protomaps/go-pmtiles v1.3.0

require (
	github.com/aaronland/go-http-bootstrap v0.1.0
	github.com/aaronland/go-http-ping/v2 v2.0.0
	github.com/aaronland/go-http-rewrite v1.0.1
	github.com/aaronland/go-http-sanitize v0.0.6
	github.com/aaronland/go-http-server v1.0.0
	github.com/aaronland/go-http-server-tsnet v0.9.0
	github.com/aaronland/go-http-tangramjs v0.1.3
	github.com/protomaps/go-pmtiles v1.5.1
	github.com/rs/cors v1.8.2
	github.com/sfomuseum/go-flags v0.10.0
	github.com/sfomuseum/go-geojsonld v1.0.0
	github.com/sfomuseum/go-http-auth v0.0.6
	github.com/sfomuseum/go-http-protomaps v0.0.14
	github.com/sfomuseum/go-http-tilezen v0.0.7
	github.com/sfomuseum/go-sfomuseum-pmtiles v1.0.3
	github.com/tidwall/gjson v1.14.3
	github.com/tilezen/go-tilepacks v0.0.0-20220729022432-5ee633f5bb6a
	github.com/whosonfirst/go-cache v0.5.2
	github.com/whosonfirst/go-reader v1.0.1
	github.com/whosonfirst/go-reader-cachereader v0.2.5
	github.com/whosonfirst/go-reader-findingaid v0.14.0
	github.com/whosonfirst/go-sanitize v0.1.0
	github.com/whosonfirst/go-whosonfirst-export/v2 v2.6.0
	github.com/whosonfirst/go-whosonfirst-feature v0.0.24
	github.com/whosonfirst/go-whosonfirst-image v0.1.0
	github.com/whosonfirst/go-whosonfirst-search v0.1.0
	github.com/whosonfirst/go-whosonfirst-spr-geojson v0.0.8
	github.com/whosonfirst/go-whosonfirst-spr/v2 v2.2.1
	github.com/whosonfirst/go-whosonfirst-svg v0.1.0
	github.com/whosonfirst/go-whosonfirst-uri v1.2.0
	github.com/whosonfirst/go-writer/v3 v3.1.0
	gocloud.dev v0.27.0
)

require (
	filippo.io/edwards25519 v1.0.0-rc.1 // indirect
	github.com/RoaringBitmap/roaring v1.2.1 // indirect
	github.com/aaronland/go-artisanal-integers v0.9.1 // indirect
	github.com/aaronland/go-aws-dynamodb v0.0.5 // indirect
	github.com/aaronland/go-aws-session v0.0.6 // indirect
	github.com/aaronland/go-brooklynintegers-api v1.2.6 // indirect
	github.com/aaronland/go-http-leaflet v0.1.0 // indirect
	github.com/aaronland/go-pool/v2 v2.0.0 // indirect
	github.com/aaronland/go-roster v1.0.0 // indirect
	github.com/aaronland/go-string v1.0.0 // indirect
	github.com/aaronland/go-uid v0.3.0 // indirect
	github.com/aaronland/go-uid-artisanal v0.0.2 // indirect
	github.com/aaronland/go-uid-proxy v0.0.2 // indirect
	github.com/aaronland/go-uid-whosonfirst v0.0.2 // indirect
	github.com/akrylysov/algnhsa v0.12.1 // indirect
	github.com/akutz/memconn v0.1.0 // indirect
	github.com/alexbrainman/sspi v0.0.0-20210105120005-909beea2cc74 // indirect
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129 // indirect
	github.com/aws/aws-lambda-go v1.13.3 // indirect
	github.com/aws/aws-sdk-go v1.44.134 // indirect
	github.com/aws/aws-sdk-go-v2 v1.16.8 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.15.15 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.12.10 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssm v1.27.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.10 // indirect
	github.com/aws/smithy-go v1.12.0 // indirect
	github.com/bits-and-blooms/bitset v1.2.0 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/coreos/go-iptables v0.6.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fxamacker/cbor/v2 v2.4.0 // indirect
	github.com/g8rswimmer/error-chain v1.0.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/godbus/dbus/v5 v5.0.6 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/google/wire v0.5.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/hdevalence/ed25519consensus v0.0.0-20220222234857-c00d1f31bab3 // indirect
	github.com/insomniacslk/dhcp v0.0.0-20211209223715-7d93572ebe8e // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/native v1.0.0 // indirect
	github.com/jsimonetti/rtnetlink v1.1.2-0.20220408201609-d380b505068b // indirect
	github.com/jtacoma/uritemplates v1.0.0 // indirect
	github.com/klauspost/compress v1.15.4 // indirect
	github.com/kortschak/wol v0.0.0-20200729010619-da482cc4850a // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	github.com/mdlayher/genetlink v1.2.0 // indirect
	github.com/mdlayher/netlink v1.6.0 // indirect
	github.com/mdlayher/sdnotify v1.0.0 // indirect
	github.com/mdlayher/socket v0.2.3 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/natefinch/atomic v1.0.1 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/paulmach/go.geojson v1.4.0 // indirect
	github.com/paulmach/orb v0.7.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/schollz/progressbar/v3 v3.11.0 // indirect
	github.com/sfomuseum/go-edtf v1.1.1 // indirect
	github.com/sfomuseum/go-tilezen v0.0.6 // indirect
	github.com/srwiley/oksvg v0.0.0-20220731023508-a61f04f16b76 // indirect
	github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef // indirect
	github.com/tailscale/certstore v0.1.1-0.20220316223106-78d6e1c49d8d // indirect
	github.com/tailscale/golang-x-crypto v0.0.0-20221009170451-62f465106986 // indirect
	github.com/tailscale/goupnp v1.0.1-0.20210804011211-c64d0f06ea05 // indirect
	github.com/tailscale/netlink v1.1.1-0.20211101221916-cabfb018fe85 // indirect
	github.com/tcnksm/go-httpstat v0.2.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	github.com/u-root/uio v0.0.0-20220204230159-dac05f7d2cb4 // indirect
	github.com/vishvananda/netlink v1.1.1-0.20211118161826-650dca95af54 // indirect
	github.com/vishvananda/netns v0.0.0-20211101163701-50045581ed74 // indirect
	github.com/whosonfirst/go-geojson-svg v0.0.5 // indirect
	github.com/whosonfirst/go-ioutil v1.0.2 // indirect
	github.com/whosonfirst/go-reader-http v0.3.1 // indirect
	github.com/whosonfirst/go-whosonfirst-findingaid/v2 v2.7.1 // indirect
	github.com/whosonfirst/go-whosonfirst-flags v0.4.4 // indirect
	github.com/whosonfirst/go-whosonfirst-format v0.4.1 // indirect
	github.com/whosonfirst/go-whosonfirst-id v1.0.0 // indirect
	github.com/whosonfirst/go-whosonfirst-placetypes v0.3.0 // indirect
	github.com/whosonfirst/go-whosonfirst-sources v0.1.0 // indirect
	github.com/whosonfirst/go-writer-featurecollection/v3 v3.0.0-20220916180959-42588e308a3e // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/ratelimit v0.2.0 // indirect
	go4.org/mem v0.0.0-20210711025021-927187094b94 // indirect
	go4.org/netipx v0.0.0-20220725152314-7e7bdc8411bf // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e // indirect
	golang.org/x/image v0.0.0-20211028202545-6944b10bf410 // indirect
	golang.org/x/net v0.1.0 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/term v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	golang.org/x/time v0.0.0-20220722155302-e5dcc9cfc0b9 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	golang.zx2c4.com/wintun v0.0.0-20211104114900-415007cec224 // indirect
	golang.zx2c4.com/wireguard v0.0.0-20220904105730-b51010ba13f0 // indirect
	golang.zx2c4.com/wireguard/windows v0.5.3 // indirect
	google.golang.org/api v0.91.0 // indirect
	google.golang.org/genproto v0.0.0-20220802133213-ce4fa296bf78 // indirect
	google.golang.org/grpc v1.48.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gvisor.dev/gvisor v0.0.0-20220817001344-846276b3dbc5 // indirect
	modernc.org/libc v1.16.7 // indirect
	modernc.org/mathutil v1.4.1 // indirect
	modernc.org/memory v1.1.1 // indirect
	modernc.org/sqlite v1.17.3 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
	tailscale.com v1.32.2 // indirect
	zombiezen.com/go/sqlite v0.10.1 // indirect
)
