#! /bin/bash

go build -o bookings cmd/web/*.go
./bookings -dbname=bookings -dbuser=corey.slate -cache=false -production=false