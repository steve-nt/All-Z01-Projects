#!/bin/bash

ls -l | awk 'NR > 1 && NR % 2 == 0'

