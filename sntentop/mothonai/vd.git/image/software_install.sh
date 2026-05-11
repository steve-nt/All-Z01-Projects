#!/bin/sh
dnf install -y $(cat /sl|tr '\n' ' ')
