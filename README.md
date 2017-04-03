# LinuxPA
LinuxPA is a try to bring a [PortableApps.com](http://portableapps.com) type launcher to Linux.  

# How to use
Just double click on an app to launch it! If there are multiple executables, you can either select the specific executable, or if you just double click the app it'll launch the first linux executable (.sh script files have priority), but if one isn't found it launches the first executable in general.

# Apps:
Both of the below places provide linux executables that don't need libs installed on the host system:  
[AppImage](https://bintray.com/probono/AppImages)  

# PortableApps.com Compatibility
LinuxPA works will with the PortableApps.com launcher, as it looks for apps in the PortableApps folder and grabs the app's name and icon from where it should be in the PortableApps.com format.  
My forum at PortableApps.com can be found [here](http://portableapps.com/node/54998).  

# common.sh
common.sh is found in the PortableApps/LinuxPACom folder and is executed before the app. I mainly use it to set environment variables (such as HOME).  

# Simple App Setup
Because apps aren't natively formated in the PortableApps.com format, if LinuxPA doesn't find the AppInfo.ini or appicon_\*.png in the App/AppInfo folder of the app it looks for them in the root directory of the app (except it looks, nor for appicon_\*.png, but appicon.png). If an AppInfo.ini file isn't found then the name of the app is grabbed from the folder name and it's category is set to other. It specifically looks for the lines starting with `Name=` and `Category=`  

# AppImage Support
[AppImage Website](http://appimage.org)  
Right now AppImages are simply supported via the native linux executable support, but later I'm hoping to add downloading and automatic updating support later on.  

# USB mount
Unfortunately Linux, by default, doesn't support running executables off of FAT formated flash drives, requiring you to mount your drive with special mount arguments or format in a linux friendly format (such as EXT4). I personally use the arguments `exec,noauto,nodev,nosuid,umask=0000`  

# Screenshots
Photos are found [Here](https://goo.gl/photos/VtBUL6DyZTMidj5n6)

# TODO (Might be in order)
1. MAKE IT BETTER  
1. Add settings menu  
1. Add updater for .AppImage files  
1. Download .AppImage files (maybe)  
1. Check if all apps are closed when it closes and ask if you want to force stop the apps.  
1. Portable wine (Should be able to get a working version from PlayOnLinux, but I need to add a check to see if filesystem is EXT as Wine doesn't like filesystems w/o permission control)
