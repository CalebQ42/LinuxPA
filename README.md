# LinuxPA
The goal is to create a fully functional PortableApps.com type launcher that can properly parse data from the PortableApps.com format. Apps are launched by a .sh file in the app's directory. Currently pulls out the Name and Category from App/AppInfo/appinfo.ini  
Works well with AppImage apps.

# Why?
I know that Linux only has about 2% desktop usage and I know that the traditional way to install apps isn't portable, but over the past year or so I've started to put linux apps on my flash drive (AppImage is a great example of a portable solution to linux apps. Also a lot of DRM-free games can be run portably), but there was no easy way to organize my linux apps, so I created one. I personally have used the PortableApps.com launcher for years now and I love how properly formated the apps are, which allows me to grab info about the app easily.  

# Why script files?
In general linux executable files have no extensions and can be a pain when trying to figure out what is executable and what isn't. I figured script files are easy to detect and allow a large amount of flexibility for me (and others who want to make apps work with this launcher). See below for .AppImage support (Get AppImages from [here](https://bintray.com/probono/AppImages))

# Why Go?
Because I like Go :) Also the way it includes all it needs into one friendly executable.

# What is needed?
Basically you need go to compile the source, AND YOU ALSO NEED TO MOUNT YOUR FLASH DRIVE SO YOU CAN EXECUTE FILES ON IT!!!! I've found that the mount arguments of `exec,noauto,nodev,nosuid,umask=0000` works well (I personally put my flash drive into /etc/fstab).

# Format
The first place the program looks for an app's icon and info is in the /App/AppInfo directory (icon defaults to appicon_32.png, otherwise it just picks the last one it finds), but if it can't find the appinfo.ini or app icon, it looks in the apps root directory for appinfo.ini and appicon.png for info and icon respectively(Just to make it easier for custom settings in an app).

# common.sh
common.sh is run before any program so you can set environment variables (such as HOME). common.sh should be in PortableApps/LinuxPACom folder.

# AppImage support
It will now launch .AppImage files! If a .sh script and an .AppImage executable are both in a directory, the .sh script takes precedence. You can get AppImages from [here](https://bintray.com/probono/AppImages).

# TODO (Might be in order)
1. MAKE IT BETTER  
1. Improve linux executable detection (A.K.A. a pain in the butt)    
1. Launching of .exe files via wine (wine will have to be installed on the host system, unless there is some portable wine, I may have found one)  
1. Add settings menu  
1. Add updater for .AppImage files   
1. Download .AppImage files (maybe)  
1. Check if all apps are closed when it closes and ask if you want to force stop the apps.  
