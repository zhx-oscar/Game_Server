module Daisy

go 1.14

replace Cinder => ../Cinder

require (
	Cinder v0.0.0-00010101000000-000000000000
	github.com/ByteArena/box2d v1.0.2
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-redis/redis/v7 v7.4.0
	github.com/gogo/protobuf v1.3.1
	github.com/jessevdk/go-flags v1.4.0
	github.com/json-iterator/go v1.1.10
	github.com/magicsea/behavior3go v0.0.0-20200622063830-4cf5449990a7
	github.com/prometheus/common v0.4.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/sipt/GoJsoner v0.0.0-20170413020122-3e1341522aa6
	github.com/spf13/viper v1.7.1
	go.mongodb.org/mongo-driver v1.4.0
)
