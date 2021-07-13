#!/bin/bash

protoc phonebook/phonebookpb/phonebookpb.proto --go_out=. --go-grpc_out=.
protoc calculator/calculatorpb/calculatorpb.proto --go_out=. --go-grpc_out=.