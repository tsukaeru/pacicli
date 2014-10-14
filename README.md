pacicli
=======
`pacicli` is a command line interface tool for Parallels Cloud Infrastructure
(PACI) API. It makes it easy to manage Containers and Virtual machines on PACI.

## Features

- Works on most of major platforms like Linux, Windows, MacOS X etc.
- Full API support described in the [Official API Document](http://download.pa.parallels.com/poa/5.5/doc/pdf/POA%20RESTful%20API%20Guide/poa_5.5_paci_restful_api_guide.pdf)
- JSON, TOML output support.

## Install

Download your platform binary from [Release page](https://github.com/tsukaeru/pacicli/releases)
and put it somewhere you like.

If you'd like to build it by yourself, please use `go get`:

```bash
$ go get -d github.com/tsukaeru/pacicli
```

## Usage

1. Retrieve your PACI API key from your Hosting Provider's web interface.
2. Create your working directory, go into it and put `Pacifile` configuration file.
   `Pacifile` should have following settings.

   ```toml
   BaseURL  = "https://example.com/paci/v1.0" # API URL
   Username = "username" # Your account name
   Password = "password" # Your API key
   ```

   You can write `Pacifile` both in JSON and TOML. `pacicli` detects its format
   automatically. If the file begins with `{` character, it's parsed as JSON.
   If not, parsed as TOML.
3. Run

   ```bash
   pacicli list
   ```

   If you have some machiens, those would be displayed. If not, please try

   ```bash
   pacicli oslist
   ```

   It would show the OS list provided by the Hoster. For more command detail,
   please see

   ```bash
   pacicli help
   ```

## Contribution

1. Fork ([https://github.com/tsukaeru/pacicli/fork](https://github.com/tsukaeru/pacicli/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `go fmt`
7. Create new Pull Request

## License

`pacicli` is under MIT license. See the [LICENSE](./LICENSE) file for details.
