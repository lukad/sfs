Simple File Server
==================

## Installation

	$ go get -u github.com/lukad/sfs

## Usage
	
	$ sfs --help
	Usage: sfs [options]
	Where options are:
	  -d, --digest=false: Use digest access authentication
	  -l, --listen=":8080": Listen address
	  --log=true: Log to stdout
	  -n, --no-color=false: Don't log with colors
	  -p, --password="": Password for authentication
	  -r, --root="./": Root directory for the file server
	  -u, --user="": Username for authentication

	$ sfs -d -ufoo -pbar -r/home/foo/baz/

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
