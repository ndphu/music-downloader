# music-downloader
This cli downloads music from multiple source (zing, nhaccuatui,...)

## Features
* Download entire album or single track
* Supproted music providers: 
    * zing (mp3.zing.vn)
    * nct (nhaccuatui.com)
* Download High Quality track (if logged with a VIP account)

## Installation
Download the prebuilt binary found or install from source
```
go get github.com/ndphu/music-downloader
go install github.com/ndphu/music-downloader
```

## Usage
### list supported providers
`music-downloader provider ls`

Output:
```
Supported providers:
  - zing
  - nct
```
### login
`music-downloader provider login --name nct -u <username> -p <password>`

Output:
```
Login...
Login successfully!
```
NOTE: Your login information will be cached in `$HOME/.music-downloader/auth`, so you don't need to login every time you use the CLI.

### download
`music-downloader download --output <output_dir> --thread-count <number_of_parrallel_download> <link_1> <link_2>...<link_n>`

If no `output` provided, it will use the current directory.
We don't use parallel download by default.

NOTE: For album, tracks are downloaded into a folder of the album name.

## Authors
Phu Nguyen <ngdacphu.khtn@gmail.com>

## TODO
* Login with zing account
