#!/bin/bash
set -eu

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
BIN_DIR="${SCRIPT_DIR}/bin/pkg"
CMD_WAVELET_DIR="${SCRIPT_DIR}/../cmd/wavelet"
VERSION=$(${BIN_DIR}/linux-amd64/wavelet -v | grep "^Version:" | awk '{print $2}')

# clean up old zip files
rm ${SCRIPT_DIR}/bin/*.zip

# loop through all the platforms and create an archive
cd ${SCRIPT_DIR}
for PLATFORM_DIR in ${BIN_DIR}/*; do
    PLATFORM=$(basename ${PLATFORM_DIR})
    echo "Archiving platform ${PLATFORM}"

    # copy over the auxilary files
    cp -R ${CMD_WAVELET_DIR}/services ${BIN_DIR}/${PLATFORM}
    rm ${BIN_DIR}/${PLATFORM}/services/README.md
    cp ${CMD_WAVELET_DIR}/config.toml ${BIN_DIR}/${PLATFORM}
    cp ${CMD_WAVELET_DIR}/genesis.json ${BIN_DIR}/${PLATFORM}
    cp ${CMD_WAVELET_DIR}/wallet.txt ${BIN_DIR}/${PLATFORM}
    cp ${CMD_WAVELET_DIR}/wallets.txt ${BIN_DIR}/${PLATFORM}

    cd ${BIN_DIR}/${PLATFORM}
    zip -r "${SCRIPT_DIR}/bin/wavelet-${VERSION}-${PLATFORM}.zip" *
done

echo "Release done"
