#!/bin/bash

cd ../v2/cmd/nuclei
go build
cp nuclei ../../../integration_tests/nuclei 
cd ../integration-test
go build
cp integration-test ../../../integration_tests/integration-test 
cd ../../../integration_tests
./integration-test
if [ $? -eq 0 ]
then
  exit 0
else
  exit 1
fi
