# LinuxPA
LinuxPA is a try to bring a [PortableApps.com](http://portableapps.com) type launcher to Linux.  

# How to use
Just double click on an app to launch it! If there are multiple executables, you can either select the specific executable, or if you just double click the app it'll launch the first linux executable it finds. .sh script files have priority over other executable files.

# Apps:
The below place provides linux executables that don't need libs installed on the host system:  
[AppImage](https://bintray.com/probono/AppImages)  

# PortableApps.com Compatibility
LinuxPA works will with the PortableApps.com launcher, as it looks for apps in the PortableApps folder and grabs the app's name and icon from where it should be in the PortableApps.com format.  
My forum at PortableApps.com can be found [here](http://portableapps.com/node/54998).  

# common.sh
common.sh is found in the PortableApps/LinuxPACom folder and is executed before the app. I mainly use it to set environment variables (such as HOME). You can create and edit the common.sh from settings  

# Simple App Setup
Because apps aren't natively formated in the PortableApps.com format, LinuxPA will look in the root directory for a AppInfo.ini (for basic info such as category and name) and appicon.png. If they aren't found, it looks where the appicon_\*.png and AppInfo.ini is in PortableApps format. You can set what the AppInfo.ini and appicon.png are from LinuxPA.  

# AppImage Support
[AppImage Website](http://appimage.org)  
Right now AppImages are simply supported via the native linux executable support, and you can download AppImages. (Woo)

# USB mount
Unfortunately Linux, by default, doesn't support running executables off of FAT formated flash drives, requiring you to mount your drive with special mount arguments or format in a linux friendly format (such as EXT4). I personally use the arguments `exec,noauto,nodev,nosuid,umask=0000`  

# Screenshots
Photos are found [Here](https://goo.gl/photos/VtBUL6DyZTMidj5n6). The screenshots are with the adapta gtk theme

# TODO (Might be in order)
1. MAKE IT BETTER  
1. Ask if you want to update  
1. Better appimage support in general   
1. Check if all apps are closed when it closes and ask if you want to force stop the apps. . 
