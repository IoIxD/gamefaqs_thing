quick and dirty (very dirty, i'm aware this code isn't good, i spent like 30 minutes on this) web server that tells you what gamefaqs pages aren't archived yet.

## Hosting this on your own machine

basic knowledge of the terminal and whatnot is assumed

- get the code through that "download zip" button up there
- [get golang](https://go.dev/)  (in the installer make sure golang is added to your PATH, if you're on Windows)
- open a terminal/command prompt/power shell
- cd to the directory where you downloaded everything.
- run `go run .`. you'll need to wait a bit while dependencies download. If you get an error about an recognized command, go is not in your PATH. do also note this is untested on windows. when all is done, you should get a "running" message
- go to "localhost:8085" in your browser.
