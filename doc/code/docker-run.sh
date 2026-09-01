docker run -d --name yarr \
  -p 7070:7070 \
  -v yarr-data:/data \
  nkanaev/yarr:latest -db /data/yarr.db -addr 0.0.0.0:7070
