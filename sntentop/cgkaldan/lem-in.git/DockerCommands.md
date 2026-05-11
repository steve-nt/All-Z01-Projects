# Docker Commands
- docker build --no-cache -t antfarm:latest .
- docker run -it --rm -v $(pwd)/example01.txt:/usr/local/bin/example01.txt   -v /tmp/.X11-unix:/tmp/.X11-unix   -e DISPLAY="$DISPLAY" antfarm example01.txt -v

### **Still not working**