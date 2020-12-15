#!/bin/bash

git add feed/*
now=$(date)
git commit -m "feat: Feed generation for [${now}]"
