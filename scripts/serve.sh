# A script to start/restart the server to listen to HTTP request on every change to go files
while sleep 1;
  do
    find *.go | entr -r go run *.go;
  done
