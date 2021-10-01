#!/bin/bash

#Special variables:
#$ROOT = directory where LinuxPA is
#$APPNAME = name of the launched app
#$FILENAME = name of the launched file

export HOME=$ROOT/PortableApps/LinuxPACom/Home
export XDG_CONFIG_HOME=$HOME/.config

#Run the app. DON'T TOUCH (unless you know what your doing :P)
$@