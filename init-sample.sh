#!/usr/bin/env bash

jq -n '{ "add.avg": 4570, "delete.avg": 353, "update.avg": 2015, "loc.total": 307, "loc.covered": 150 }' | http -a "$AUTH" POST ':8000/project/ff/report/0.0.1'
jq -n '{ "add.avg": 4750, "delete.avg": 335, "update.avg": 2035, "loc.total": 505, "loc.covered": 213 }' | http -a "$AUTH" POST ':8000/project/ff/report/0.0.2'
jq -n '{ "add.avg": 8530, "delete.avg": 335, "update.avg": 2045, "loc.total": 750, "loc.covered": 213 }' | http -a "$AUTH" POST ':8000/project/ff/report/0.0.4'
jq -n '{ "add.avg": 4775, "delete.avg": 335, "update.avg": 2035, "loc.total": 730, "loc.covered": 471 }' | http -a "$AUTH" POST ':8000/project/ff/report/0.0.7'
jq -n '{ "add.avg": 4570, "delete.avg": 335, "update.avg": 2005, "loc.total": 800, "loc.covered": 511 }' | http -a "$AUTH" POST ':8000/project/ff/report/1.0.0'
jq -n '{ "add.avg": 4120, "delete.avg": 450, "update.avg": 2035, "loc.total": 854, "loc.covered": 530 }' | http -a "$AUTH" POST ':8000/project/ff/report/1.2.1'
jq -n '{ "add.avg": 4220, "delete.avg": 335, "update.avg": 2035, "loc.total": 954, "loc.covered": 533 }' | http -a "$AUTH" POST ':8000/project/ff/report/1.3.1'
