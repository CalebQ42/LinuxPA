# LinuxPA
LinuxPA is a try to bring a [PortableApps.com](http://portableapps.com) type launcher to Linux.  

# App Detection
LinuxPA looks in all folders in the PortableApps folder for, first, a script file (starts with the shebang (`#!`)) and, secondly, a native linux executable (starts with ELF). It will only add the first one it finds.  

# PortableApps.com Compatibility
LinuxPA works will with the PortableApps.com launcher, as it looks for apps in the PortableApps folder and grabs the app's name and icon from where it should be in the PortableApps.com format.  

# common.sh
common.sh is found in the PortableApps/LinuxPACom folder and is executed before the app. I mainly use it to set environment variables (such as HOME).  

# Simple App Setup
Because apps aren't natively formated in the PortableApps.com format, if LinuxPA doesn't find the AppInfo.ini or appicon_\*.png in the App/AppInfo folder of the app it looks for them in the root directory of the app (except it looks, nor for appicon_\*.png, but appicon.png). If an AppInfo.ini file isn't found then the name of the app is grabbed from the folder name and it's category is set to other. It specifically looks for the lines starting with `Name=` and `Category=`  

# AppImage Support
[AppImage Website](http://appimage.org)  
Right now AppImages are simply supported via the native linux executable support, but later I'm hoping to add downloading and automatic updating support.  

# USB mount
Unfortunately Linux, by default, doesn't support running executables off of flash drives, requiring you to mount your drive with special mount arguments, I personally use the arguments `exec,noauto,nodev,nosuid,umask=0000`  

# Screenshots
Photos are found [Here](https://goo.gl/photos/VtBUL6DyZTMidj5n6)

# TODO (Might be in order)
1. MAKE IT BETTER  
1. Add settings menu  
1. Add updater for .AppImage files  
1. Download .AppImage files (maybe)  
1. Check if all apps are closed when it closes and ask if you want to force stop the apps.  
