#!/bin/bash -e

TMP_DIR="../tmp"
TARGET_REPO="https://github.com/golang/tools.git"

if [ -z "$1" ] ;then
    echo "Target repo supplied: will use $1"
    TARGET_REPO=$1
else
    echo "Target repo was not supplied: will use default $TARGET_REPO"
fi

mkdir -p "$TMP_DIR"
pushd "$TMP_DIR"

echo "Getting latest TruffleGopher.."
go get "github.com/michaelgrifalconi/trufflegopher/cmd/cli" &> /dev/null

echo "Getting latest TruffleHog.."
pip install truffleHog &> /dev/null


echo "Writing sample TruffleHog Rule"
cat << EOF > trufflehog-rules.json
{
    "AWS KEY": "AKIA[0-9A-Z]{16}"
}
EOF

echo "Writing sample TruffleGopher Rule"
cat << EOF > trufflegopher-rules.json
signatures:
- exp:  "AKIA[0-9A-Z]{16}"
  description:  "AWS KEY"
EOF

echo "Cloning benchmark target repository"
git clone --quiet "$TARGET_REPO" target-repo
TARGET=$(pwd)/target-repo

echo "STARTING Trufflehog scan.."
(time truffleHog --rules trufflehog-rules.json file://"$TARGET" ) > /dev/null

echo "STARTING Trufflegopher scan.."
(time ./tg-cli -signatures="trufflegopher-rules.json" -repo="$TARGET")

popd

echo "Cleaning up temp dir"
rm -rf $TMP_DIR