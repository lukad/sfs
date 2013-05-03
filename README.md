Simple File Server
==================

## Installation

	$ go get -u github.com/lukad/sfs

## Usage
	
	$ sfs --help
	Usage: sfs [options]
	Where options are:
	  -d, --digest=false: Use http digest authentication
	  -l, --listen=":8080": Listen address
	  -p, --password="": Password for basic auth
	  -r, --root="./": Root Directory for the file server
	  -u, --user="": Username for basic auth

	$ sfs -d -ufoo -pbar -r/home/foo/baz/

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
