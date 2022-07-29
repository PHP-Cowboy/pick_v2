# git pull

go build -o pick_v2 main.go

killall pick_v2

nohup ./pick_v2 >> nohup.log 2>&1 &