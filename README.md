# LinuxPA

LinuxPA is a try to bring a [PortableApps.com](http://portableapps.com) type launcher to Linux.

(This branch is for UI testing and things might change faster then the README)

## How to use

Just double click on an app to launch it! If there are multiple executables, you can either select the specific executable, or if you just double click the app it tries to launch the best executable. Order of priority is Desktop Files > scripts > AppImages > native linux executables > windows apps (if wine is enabled).

## Apps

The below place provides linux executables that don't need libs installed on the host system:

* [https://appimage.github.io/](https://appimage.github.io/)

## PortableApps.com Compatibility

LinuxPA works will with the PortableApps.com launcher, as it looks for apps in the PortableApps folder and grabs the app's name and icon from where it should be in the PortableApps.com format.

## common.sh

common.sh is a script file that gets executed at launch. This is mainly to set thing such as HOME for further portability. common.sh is enabled by default, but can be disabled per app or globally.

NOTE: This behavior is different from v2 where common.sh was disabled by default. Additionally, instead of being launched seperately, programs are launched from the script. This should help with compatibility, but if your script file is borked, the app will not launch.

## Simple App Setup

Because apps aren't natively formated in the PortableApps.com format, LinuxPA will look in the root directory for a AppInfo.ini (for basic info such as category and name) and appicon.png. If they aren't found, it looks where the appicon_\*.png and AppInfo.ini is in PortableApps format. You can set what the AppInfo.ini and appicon.png are from LinuxPA.

## AppImage Support

[AppImage Website](http://appimage.org)
Some advanced features are available for AppImages including:

* Icons and Desktop files are obtained and parsed from within the AppImage.
* Create .config and .home folders for AppImages on launch

## USB mount

LinuxPA, by default, WILL NOT work on non-unix file systems (such as NTFS and FAT32). If you wish to keep your flash drive formatted as a non-unix file system (probably to keep compatibility with Windows and the PortableApps.com Launcher) you have to mount the flash drive with the options `exec,noauto,nodev,nosuid,umask=0000`.

## Screenshots

Photos are found [Here](https://goo.gl/photos/VtBUL6DyZTMidj5n6). The screenshots are with the adapta gtk theme

## TODO (Might be in order)

1. MAKE IT BETTER
1. Manual update check
1. Better AppImage integrations (Specifically updating and appimage downloading from LinuxPA)
1. Sandboxing support
   1. Might be possible by packaging as an AppImage and providing Firejail, or simply just downloading it (like with Wine)
1. Check if all apps are closed when it closes and ask if you want to force stop the apps
1. Window showing the output of the apps after they're launched.
