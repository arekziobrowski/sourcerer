# sourcerer
![test workflow](https://github.com/arekziobrowski/sourcerer/actions/workflows/test.yml/badge.svg)


Git-versioned source code and Maven dependency download.

# Source code download
Currently, only Git-versioned source code download is supported.

The following steps are done in order to ensure the best performance (both time and space-wise):

```shell
git init
git remote add origin git@github.com:<ORGANIZATION>/<REPO>.git
git fetch origin <SHA1> --depth=1
git reset --hard FETCH_HEAD
```
