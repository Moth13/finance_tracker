# finance_tracker
A personal finance tracker backend and UI, using postgresql, golang, gin, htmx and css.

## Installation

Make sure you have installed golang and docker.

```
git clone https://github.com/Moth13/finance_tracker.git
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

