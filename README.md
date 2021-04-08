# sourcerer
![test workflow](https://github.com/arekziobrowski/sourcerer/actions/workflows/test.yml/badge.svg)


Git-versioned source code download.

# Source code download
Currently, only Git-versioned source code download is supported.

The following steps are done in order to ensure the best performance (both time and space-wise):

```
git init
git remote add origin git@github.com:<ORGANIZATION>/<REPO>.git
git fetch origin <SHA1> --depth=1
git reset --hard FETCH_HEAD
```

# Input file format
Input list should be structured in the following manner:
```shell
git@github.com:org/repo-name-1.git <revision hash>
git@github.com:org/repo-name-2.git <revision hash>
```
# Dependency target directories
Maven dependencies are downloaded based on `.sourcerer-pom.xml` file in the project directory. The downloaded dependency jars are placed in `.sourcerer-deps` directory.
