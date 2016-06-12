# Pushabutton
Minimal web app to run scripts.

Put anything executable in the `buttons` dir.
It'll show up as a button in the web UI.
Then, you can push a button.

For example,

    buttons$ ln -s $(which date)

Looking at the web-UI, you immediately get a button that you can click.



![Pushabutton automatically renders a button for you](docs/intro-date-button.png)



After clicking, the executable is run and you get a log:


![Pushabutton shows you what is happening on pushing the button](docs/intro-date-log.png)


Remember what command you ran yesterday? No worries, `Pushabutton` remembers.


![Pushabutton remembers your button pushes](docs/intro-date-log-list.png)

## Simple Parameters

Oh, and then there's parameters! Let's say you're interested in running the `id` command from the web, to find the IDs of users:

    $ ln -s $(which id)

Now, you need to pass a username. Add a textfile `id.parameters`:

    $ cat id.parameters
    username

When clicking the button, you'll get a modal asking for the username:

![Pushabutton automatically asks for the username!](docs/intro-id-parameters-modal.png)

## More features

* Add descriptions to parameters
* Make the user choose a value for a parameter with a drop-down list.
* Determine the valid choice at runtime (by running a script).

How to configure these, see the examples after running setup.


## What it is and what it's not

Pushabutton is about enabling non-techies to use your simple helper scripts, too.
You copy the binary, you dump the script in the buttons dir, you run the webserver. Then you paste the link in the mail to your colleague, and that's it — done.
You can do simple parameters, you have a simple log — but that's it.
Pushabutton fills a gap at the low-effort, simple-need end of the spectrum.

It will not be enough for many use cases, and in fact it is not meant to be. If you need:

* Authorization (this user can do this, the other one cannot)
* Scheduling (at 5pm tomorrow do this and at every 2nd Sunday do that)
* Editing from the web
* Pipelining
* Conditional logic

then Pushabutton is not for you.
There are other tools that where you can run something from the web, and that do those things well:

* [Rundeck](http://rundeck.org/),
* [Jenkins](https://jenkins.io/index.html)
* and many others.


## Setup & Running & Uninstalling

Run the setup command:

  $ ./pushabutton setup

This will setup the necessary directories and some demo scripts:

    $ tree
    .
    ├── buttons
    │   ├── what_is_the_current_date.sh
    │   ├── write-message-to-logged-in-user.sh
    │   ├── write-message-to-logged-in-user.sh.parameters
    │   └── write-message-to-logged-in-user.sh.parameters.user-tty.choices.sh
    ├── logs
    └── pushabutton

What did we get?

* `buttons` is where your scripts live. Link them or drop them there — as long as it's executable, Pushabutton will make a button out of it. At the beginning, two scripts demonstrate how to create your own buttons.
* `logs` that's where, for every button push, log files are stored. The `journal.log` is the file that indexes the rest.


**Run the server**:

   $ ./pushabutton serve
   Running server on :8080


Want to **uninstall**? Simple. Just delete the entire directory.
There's nothing else.

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
