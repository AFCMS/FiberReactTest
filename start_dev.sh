#!/bin/bash

FiberReactTest &

cd frontend && npm run start &

wait -n

exit $?