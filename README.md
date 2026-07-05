# Code-the-Classics vol.1: Golang ports

The book [Code the Classics, vol.1](https://magazine.raspberrypi.com/books/code-the-classics-vol-I-2ed) describes several games programmed using Python and the pygame-zero library. The original codes and assets are provided at https://github.com/raspberrypipress/Code-the-Classics-Vol1

This fork contains ports of the original Python codes to the [Go programming language](http://go.dev), using the [go-sdl3](https://github.com/Zyko0/go-sdl3). 
The Go code was generate by Claude code under my supervision. Note: we had to move the assets in the go subfolders to embed the asses in the go binaries, so the python script will no longer find the assets. Use the original repo if you want to run the python scripts.

One of the advantages of Go with respect to Python is that Go can produce self-contained binaries which can be redistributed.

(TODO: implement CD/CI to release binaries for the most common platform on github.)


Christophe Pallier <christophe@pallier.org>  2026-07-05



