## Installation

You need to have docker installed. The version used is:

```bash
$ docker --version
Docker version 18.03.1-ce, build 9ee9f40
```

## Image 

The docker image `alextanhongpin/scraper:1.0.0` has been pushed to Dockerhub. To build your own image:

```bash
$ docker build -t orgname/scraper .
```

Alternatively, use the `Makefile` provided:

```
$ make docker
```

## Preview Built Image

The image is built using docker-multi-stage build for golang. Image size is kept to the minimum:

```bash
$ docker images | grep scraper
```

Output:

```
alextanhongpin/scraper                          latest              6040952e7cca        3 hours ago         11MB
alextanhongpin/scraper                          1.0.0               6ce1f2aef255        3 hours ago         10.9MB
```

## Run

```bash
$ docker run -it -v "$PWD:/data" alextanhongpin/scraper -o /data/repos.csv

# or

$ make run-docker
```

Output:

```
Enter the repository name, e.g. kubernetes/charts:
kubernetes/charts
Enter repository name, or [Y] to proceed:
y
Fetching 1 repositories...
updated_at,name,login,clone_url
2018-05-24T14:26:26Z,charts,kubernetes,https://github.com/kubernetes/charts.git
Done. Press ctrl + c to cancel.
```

CSV output:

```csv
updated_at,name,login,clone_url
2018-05-24T14:26:26Z,charts,kubernetes,https://github.com/kubernetes/charts.git
```