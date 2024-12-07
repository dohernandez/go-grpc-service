#!/usr/bin/env bash

# Init script to kick-start your project
url=$(git remote get-url origin)

url_nopro=${url#*//}
url_noatsign=${url_nopro#*@}

gh_repo=${url_noatsign/"github.com:"/"github.com/"}
gh_repo=${gh_repo%".git"}