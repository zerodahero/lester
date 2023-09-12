# Lester

The Listening Tester

## Installation

Download the latest `lester` and plop that thing in your environment (e.g. inside docker) at the root directory

## Usage

Once `lester` is in your root, run it with one of the following modes:

### Configuration

You'll need a configuration file, `lester.yml` to define some locations. Lester is capable of testing multiple components at the same time (i.e. if you have packages or subdependencies with separate test suites, Lester will handle those and run tests for them as well.)

```yml
testConfig:
    components:
        root: .
        # other components can go here
        # for example:
        # library: path/to/library
    aliases: # optional
        # you can also setup aliases in case your
        # component names are too boring to type
        # for example:
        # lib: library
```

### Test

```bash
./lester test <component>
```

`lester test` will run the phpunit tests for any given component. This has the added benefit of being able to do this without changing directories, but is otherwise the same as existing scripts.

#### Examples

```bash
# Test all of web
./lester test web

# Test the "banana" group in core
./lester test core -- --group banana

# Test a specific test file in lib
./lester test lib -- /path/to/tests/unit/SomeTest.php
```

### Watch

Performs the same tests as above, but watches for modifications to any files in that component and re-runs the tests after files are changed. `watch` starts with a test run.

#### Examples

```bash
# Test and watch all of web
./lester watch web

# Test and watch the "banana" group in root
# This isn't too smart--it'll re-run tests even if the file you touch isn't
# related to the banana group at all
./lester watch root -- --group banana

# Test a specific test file in lib
./lester watch lib -- /path/to/tests/unit/SomeTest.php
```

### Auto

This is where the magic _really_ happens! `auto` will watch **ALL THE FILES** in project (well, all the files that matter), and will run _the specific tests_ related to those changes.

For example, if you modify a file, `SomeClass.php`, `lester` will search for the matching test in the same component--`SomeClassTest.php`. If it finds it, we'll test that file, but also look inside to see what groups are being tested for it, and _also_ test those.

`auto` keeps a running list of all the files you've modified, so as you go around modifying things your test runs will get bigger. `lester` will also run the tests for _each_ component involved--so if you edit files in root and core, both test suites will be run for the relevant files. Groups it finds will be run across _every_ component that is being tested.

If you want to bootstrap/seed `lester auto` with your current changes in your branch, you can run:

```bash
./lester auto --seed-from-git
```

Lester will attempt to get all the committed, untracked, staged, and unstaged changes from `git` (NOTE: some earlier versions of `git`, for staged and unstaged changes can take a while to find. This should only be slow the first time though, and should be much faster after that).

#### Examples

```bash
# Run lester auto
./lester auto

# Run lester auto, seed from git
./lester auto --seed-from-git
```
