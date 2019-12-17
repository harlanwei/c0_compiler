#!/bin/bash

CURRENT_PATH=$(pwd)
cd ~
if [[ -d go ]]; then
    mkdir go; mkdir go/src;
done
rm -rf ~/go/src/c0_compiler
cp -R CURRENT_PATH

