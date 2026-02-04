#!/bin/bash
date >> $1
cloc . --exclude-dir=node_modules,dist,build,.svelte-kit >> $1
tail $1
