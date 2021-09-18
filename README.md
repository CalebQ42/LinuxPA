# LinuxPA

LinuxPA is a try to bring a [PortableApps.com](http://portableapps.com) type launcher to Linux.  

## How to use

Just double click on an app to launch it! If there are multiple executables, you can either select the specific executable, or if you just double click the app it'll launch the first linux executable it finds. .sh script files have priority over other executable files.

## Apps

The below place provides linux executables that don't need libs installed on the host system:  
[https://appimage.github.io/](https://appimage.github.io/)  

## PortableApps.com Compatibility

LinuxPA works will with the PortableApps.com launcher, as it looks for apps in the PortableApps folder and grabs the app's name and icon from where it should be in the PortableApps.com format.  
My forum at PortableApps.com can be found [here](http://portableapps.com/node/54998).  

## common.sh

common.sh is found in the PortableApps/LinuxPACom folder and is executed before the app. I mainly use it to set environment variables (such as HOME). You can create and edit the common.sh from settings.

## Simple App Setup

Because apps aren't natively formated in the PortableApps.com format, LinuxPA will look in the root directory for a AppInfo.ini (for basic info such as category and name) and appicon.png. If they aren't found, it looks where the appicon_\*.png and AppInfo.ini is in PortableApps format. You can set what the AppInfo.ini and appicon.png are from LinuxPA.  

## AppImage Support

[AppImage Website](http://appimage.org)
I'm looking into improving AppImage support. As of 2.1.5.0 IF `unsquashfs` is in $PATH then some advanced AppImage support is available and it will automagically get the name and possibly the icon of it. I'm looking into better support, but it might be a while.

## USB mount

Unfortunately Linux, by default, doesn't support running executables off of FAT formated flash drives, requiring you to mount your drive with special mount arguments or format in a linux friendly format (such as EXT4). I personally use the arguments `exec,noauto,nodev,nosuid,umask=0000`  

## Screenshots

Photos are found [Here](https://goo.gl/photos/VtBUL6DyZTMidj5n6). The screenshots are with the adapta gtk theme

## TODO (Might be in order)

1. MAKE IT BETTER
1. Integrate [goappimage](https://github.com/probonod/go-appimage) library for better AppImage integration.
1. Try to `chmod +x` executables if they don't have the permission
1. Manual update check
1. Better AppImage integrations (Specifically updating and better appimage downloading)
1. Get information (such as name and icon) directly from an appimage
1. Better appimage downloading (probably based around [AppImageHub](https://appimage.github.io/apps/))
1. Sandboxing support
   1. Might be possible by packaging as an AppImage and providing Firejail, or simply just downloading it (like with Wine)
1. Check if all apps are closed when it closes and ask if you want to force stop the apps
