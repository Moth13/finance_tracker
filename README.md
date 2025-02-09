# finance_tracker
A personal finance tracker backend and UI, using postgresql, golang, gin, htmx and css.

This is a show off project, still in development.

It shows a list of transactions, associated to account, category, month, having a screen to show current and future state.

Goal is to make it deploy on the cloud.

## Installation

Make sure you have installed golang.

```
git clone https://github.com/Moth13/finance_tracker.git
```

Install pre-commit using
```
brew install pre-commit && pre-commit install
```

Install commitling using
```
npm install -g @commitlint/cli @commitlint/config-conventional
```

## Local launch

First launch the postgresql docker image, and create the db

````
make postgres && make createdb
````

Then migrate the db
```
make migrateup
```

As some forms aren't still not implemented (for account, month,...), you need to fake some datas using:
```
make faking
```

If you want to launch it normally
```
make server
```
or if you want to launch it with an autoreload
```
make air
```

## Misc commands

Dumping the db using
```
make dropdb
```

Launching the test suite
```
make test
```
