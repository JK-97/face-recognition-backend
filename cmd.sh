
# prefix="FACE-RECOGNITION-BACKEND"
#export JX_EDGEBOX_APPID="5c077dc2d1a78054b4ddc298"
#export JX_EDGEBOX_APIGATEWAY="http://192.168.6.3:8008"
# export JX_EDGEBOX_AI_FACERECOGNITION=""
# export JX_EDGEBOX_INFRA_VIDEO_CAPTURE=""

#./bin/face-recognition-backend serve --port=8080 --camera-addr=http://192.168.6.3:8088 --face-ai-addr=http://192.168.6.3:8099 --db-addr=mongodb://192.168.3.33
./bin/face-recognition-backend serve --port=18080 --camera-addr=http://192.168.0.158:8088 --face-ai-addr=http://192.168.0.158:8098 --db-addr=mongodb://192.168.0.158
