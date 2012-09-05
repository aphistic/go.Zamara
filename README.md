go.Zamara
=========
go.Zamara is a Go (http://www.golang.org) sibling of the Zamara
project (http://www.github.com/aphistic/zamara).  It is intended
as a way for me to learn Go on a project I find interesting. Due
to this fact, the libraries, naming, interfaces or ANYTHING may
change at my whim.  It's my intention to work on this library
and create it the "right way" and that will probably change over
the course of development.  If you see anything that could be 
written better or a feature you'd like added, let me know! :)

Packages
--------
- zamara: This package contains code for the command-line utility for
	  interacting with MPQ archives and files based on MPQ archives.
- mpq:	  This package contains code for reading MPQ files.
- sc2: 	  This package contains code for interpreting StarCraft II
	  game replays.

What's a Zamara?
----------------
This library started out as a C++ "port" of the Tassadar ruby library 
(see credits below) and when I was trying to think of a name of the 
project I wanted to keep the name in the Starcraft universe, hopefully 
a Protoss as well.  As I was reading through the [SC wiki] I found a 
Protoss named Zamara that sounded interesting.  She's a Protoss Preserver, 
someone who has the ability to preserve the memories of all members of 
their species.  I thought it was the perfect name for a library meant 
to read all the memories of past SC2 games!

Goals
-----
* Open Source, cross-platform Starcraft II replay parser
* Support for multiple Starcraft II replay versions
* Few external dependencies
* Provide a Go library for working with MPQ files
* Create a command-line utility for working with MPQ files
  and Starcraft II replays.
* Create a web service to provide the same functionality in a web format

Requirements
------------
go.Zamara was tested with Go version 1.0.2 on:
* OS X 10.8.1
* Ubuntu 11.10
* Windows 7

There are currently no external dependencies aside from the standard
Go library and the Zamara mpq, sc2 and zamara packages.

Building/Testing
----------------

go.Zamara uses the standard "go" command:


### To build:

	go build

### To run tests:

	go test

Contact
-------
Website: https://github.com/aphistic/go.zamara

Credits
-------
* [Tassadar](https://github.com/agoragames/tassadar) -
The initial inspiration for this project and a source of structure and processes for reading replay files.
* libmpq -
A lot of help when writing the code for reading MPQ files.
* [mpyq](https://github.com/arkx/mpyq) -
This utility was very helpful when validating the values I was getting from the MPQ when writing the MPQ code.
* [sc2reader](https://github.com/GraylinKim/sc2reader) -
The wiki is an invaluable source of information on the format of the files within the MPQ.
* [WARP](http://trac.erichseifert.de/warp) -
Another good source for how to process SC2 replay files.
* [sc2gears](https://sites.google.com/site/sc2gears/) -
A great way to validate that my SC2 data is being read correctly.

License
-------

Copyright (c) 2012, Erik Davidson
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
