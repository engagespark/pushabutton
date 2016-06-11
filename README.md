# pushabutton
Simple web app to run scripts


# License

MIT, see [LICENSE](./LICENSE) file

# Dependencies, their license & thanks!

We're using the following libraries directly:

[gouuid](https://github.com/nu7hatch/gouuid) and [kingpin.v2](gopkg.in/alecthomas/kingpin.v2)

Thru the vendor directory, we're redistributing source code of the following projects:

    $ govend -tv

* [github.com/nu7hatch/gouuid](https://github.com/nu7hatch/gouuid) MIT
* [gopkg.in/alecthomas/kingpin.v2](gopkg.in/alecthomas/kingpin.v2) MIT
  * [github.com/alecthomas/units](https://github.com/alecthomas/units) MIT
    * [github.com/stretchr/testify/assert](https://github.com/stretchr/testify) MIT
      * [github.com/pmezard/go-difflib/difflib](https://github.com/pmezard/go-difflib) MIT
      * [github.com/davecgh/go-spew/spew](https://github.com/davecgh/go-spew) ISC
        * [github.com/davecgh/go-spew/spew/testdata](https://github.com/davecgh/go-spew) ISC
  * [github.com/alecthomas/template](https://github.com/alecthomas/template) GO license
    * [github.com/alecthomas/template/parse](https://github.com/alecthomas/template) GO license
  * [github.com/alecthomas/assert](https://github.com/alecthomas/assert) MIT
    * [github.com/sergi/go-diff/diffmatchpatch](https://github.com/sergi/go-diff) MIT
      * [github.com/stretchrcom/testify/assert](https://github.com/stretchrcom/testify) MIT
    * [github.com/alecthomas/repr](https://github.com/alecthomas/repr) MIT
    * [github.com/alecthomas/colour](https://github.com/alecthomas/colour) MIT
      * [github.com/mattn/go-isatty](https://github.com/mattn/go-isatty) MIT
        * [golang.org/x/sys/unix](https://godoc.org/golang.org/x/sys/unix) Go license
