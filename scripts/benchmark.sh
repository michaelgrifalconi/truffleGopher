#!/bin/bash -e

SCRIPT_DIR=$(dirname "$0")
TMP_DIR="$SCRIPT_DIR/../tmp"
TARGET_REPO="https://github.com/golang/tools.git"

if [ $# -eq 0 ] ;then
    echo "Target repo was not supplied: will use default $TARGET_REPO"
else
    echo "Target repo supplied: will use $1"
    TARGET_REPO=$1
fi

mkdir -p "$TMP_DIR"
pushd "$TMP_DIR" > /dev/null

echo "Getting latest TruffleGopher.."
go get "github.com/michaelgrifalconi/trufflegopher/cmd/tg-cli"

echo "Getting latest TruffleHog.."
pip install truffleHog &> /dev/null


echo "Writing sample TruffleHog Rule.."
cat << EOF > trufflehog-rules.json
{
    "AWS KEY": "AKIA[0-9A-Z]{16}"
}
EOF

echo "Writing sample TruffleGopher Rule.."
cat << EOF > trufflegopher-rules.yml
signatures:
- exp:  "AKIA[0-9A-Z]{16}"
  description:  "AWS KEY"
EOF

echo "Cloning or pulling benchmark target repository.."
if [ ! -d target-repo ] ; then
   git clone --quiet "$TARGET_REPO" target-repo
else
    pushd target-repo > /dev/null
    git pull --quiet "$TARGET_REPO"
    popd  > /dev/null
fi

TARGET=$(pwd)/target-repo

echo "============================="
echo "STARTING truffleHog scan.."
START=$(date +%s)
set +e
truffleHog --entropy=False --rules trufflehog-rules.json file://"$TARGET" > /dev/null
set -e
END=$(date +%s)
TH_TIME=$((END - START))


echo "============================="
echo "STARTING truffleGopher scan.."
START=$(date +%s)
tg-cli -signatures="trufflegopher-rules.yml" -repo="$TARGET"
END=$(date +%s)
END=$(date +%s)
TG_TIME=$((END - START))
popd

echo "truffleHog:    ${TH_TIME}s"
echo "truffleGopher: ${TG_TIME}s"
## Not cleaning up, allows to run same benchmark multiple times
#echo "Cleaning up temp dir"
#rm -rf $TMP_DIR