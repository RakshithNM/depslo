while sleep 1;
  do
    find *.go | entr -r go run *.go;
  done
