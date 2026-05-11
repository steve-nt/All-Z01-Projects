#!/bin/sh
docker build -t forum-app .
docker run -p 8080:8080 forum-app