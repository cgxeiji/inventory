# Inventory

[![Version](https://img.shields.io/github/v/tag/cgxeiji/inventory?sort=semver)](https://github.com/cgxeiji/inventory/releases)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/cgxeiji/inventory)](https://pkg.go.dev/github.com/cgxeiji/inventory)
[![License](https://img.shields.io/github/license/cgxeiji/inventory)](https://github.com/cgxeiji/inventory/blob/master/LICENSE)
![Go version](https://img.shields.io/github/go-mod/go-version/cgxeiji/inventory)

Inventory is an inventory manager written in Go and served as a web page. It
was developed to keep track of items in a small scale, like research
laboratories or homes.

## How does it work?

By default, inventory will run a web server at `localhost:8080` and create a
directory at `$HOME/.inventory` to store the information of each item.

You can change the port and directory by passing:
```
$ inventory -p 5005 -d /path/to/inventory/dir
```

`localhost:8080` shows a table with the list of items in the inventory. You can
press `+ Add item` to add a new item to the inventory using the web interface
or manually create an entry at `$HOME/.inventory`.

#### Items

Each item is added to its own directory with a unique id name and consist of a
`info.yaml` metadata file, `picture.jpg` with a picture of the item, and
`location.jpg` with a photo of the return location.

Once the item is added, it will appear inside the table at the root url.
Clicking on the thumbnail of the item will open a new page with a QR code
linking that specific item. You can print this QR code and physically attach it
to the item.

#### Check in/out with QR codes

Scanning the QR code will open an `update item` interface where the user is
prompted to write their name. Clicking on `Use item!` will change the state of
the item to `in use = true` and show the name of the person using the item.

After finishing using the item, scan the QR code again to return the item. This
time, the user is prompted to upload a photo of the place the item was returned
at. Pressing on `Return item!` will change the state of the item to `in use =
false` and the location of the item can be reviewed by clicking on the
`returned` link from the main inventory url.

### Security

If only restricted users should have access to the inventory, then create a
`.htdigest` file or `.htpasswd` file and add users to the realm `inventory`.

Then, run inventory as:
```
$ inventory -c /path/to/.htdigest
```

This will block the access to the inventory main page, but the QR scan can
still be accessed by anyone.

## Installation

You can download the [release
binaries](https://github.com/cgxeiji/inventory/releases) or compile it from
source by running:
```
$ go get -v github.com/cgxeiji/inventory/inventory
```
