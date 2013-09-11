# Trusted cloud broker #

The Trusted Cloud Broker (TCB) enables storing data securely in a variety of clouds, by keeping the encryption keys separated from the data storage.

*Primary use case: provide an easy way to store data safe from prying eyes in the cloud with a minimal local setup*

Think of it as a /key-file/ store, to a cloud, secured.

Why written in Go? Well, it basically runs on any platform, from ARM to x64, all major OSes.

##TCB 101##

_* store data to TCB -> TCB encrypts, compresses, mangles filename -> TCB stores keys and metadata locally in store -> TCB uploads data to cloud of choice.*_

So, you would run a metadata store locally. Net result: data stored in public clouds, not accessible unless people come knocking at YOUR door, requiring the metadata.
An additional benefit is that querying for just metadata doesn't require hitting the object store, with all latencies etc.

The interface is a simple web server with a REST interface, see URLS specified below.

## Quickstart ##

Install Go.

go get github.com/mbk/tcb/

Copy the tcb-sample.ini to tcb.ini

go build github.com/mbk/tcb/

./tcb (from the direcory where it was built by Go, depending on your environment)

Data will be stored in /tmp, metadata store is memory backed

## Cloud back ends ##
Backends that are supported out of the box:
* local file storage (mostly for testing)
* Amazon S3
* Openstack Swift (tested against Rackspace UK)

The metadata stores are pluggable as well:
* memory backed (for quick tests)
* file-based (for single node deployments)
* RIAK backed - for stateless, multiple node deployments

## Command line flags ##

-usessl Runs with https.

-port <some_number> Runs on another port than the default 8080

-config <config_file_location> The location and name of the config file, defaults to ./tcb.ini

## REST URLs

Here are the URLs to post to, and you will see the various shell scripts testing this against an instance running on localhost.
The few shell scripts will give an idea of how you can do, and what.

* /exists/some/path/to/a/file   for HTTP verb HEAD. Checks whether some/path/to/a/file exists
* /download/some/path/to/a/file   for HTTP verb GET. Downloads some/path/to/a/file if it exists
* /upload/some/path/to/a/file   for HTTP verb PUT, POST. Uploads some/path/to/a/file to the backend configured as default
* /upload/some/path/to/a/file/to/{backend}   for HTTP verb PUT, POST. Uploads some/path/to/a/file to the backend specified in {backend}; currently local, s3, or swift are valid values. This gives the option to be specific. Note that downloads magically will fetch the file from the backend where it was stored.
* /delete/some/path/to/a/file   for HTTP verb DELETE. Deletes some/path/to/a/file from its backend
	
Then there is the option to addd key/value pairs to stored data:

* /metadata/some/path/key/a/value/b for HTTP verb PUT: Sets the value of key "a" to "b" for some/path
* /metadata/some/path/key/a/ for HTTP verb GET: Gets the value of key "a" for some/path
* /metadata/some/path/key/a/ for HTTP verb DELETE: Deletes the value of key "a" some/path

## So how is the data secured?
- every path gets replaced by a global unique id (GUID)
- every file gets compressed locally, and encrypted before uploaded to the cloud backend
- the encryption keys stay "on-site", i.e. on your local server, not based in the cloud

So in the cloud there is just a data "BLOB" with a name like aaa-111-2222-sde32

##Roadmap
- more backends (Google data store)
- path support so we can think of ...
- FUSE interface
- maybe basic authentication
- native command line tools other than curl
- sync support (eventually)

License: BSDv3.

Funding: welcome.



