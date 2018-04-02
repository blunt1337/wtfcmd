Download the latest `wtf.exe` on https://github.com/blunt1337/wtfcmd/releases/latest.
Then move that file to another location like `C:\Program Files\wtf\wtf.exe`.
You need to add this location in your $PATH, to do so, execute the following command with your location: `setx path "%path%;C:\Program Files\wtf"`.
To update the wtf command, you can just replace the file.

Now test it by opening a new powershell, and type `wtf`.
To install the autocomplete, run `wtf --autocomplete install`.
You can rename `wtf` by another name, it will still work :3