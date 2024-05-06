#!/bin/bash

go build -o go-bookings cmd/web/*.go
./go-bookings -dbuser=system -dbname=go-bookings -production=false -cache=false -dbport=5431