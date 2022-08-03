#!/bin/bash

read -r -d '' header << EOF
// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

EOF

for file in "$@"
do
    if [[ -z "$(sed -ne '1{/^\/\/ This file is subject to the terms and conditions/!q;}' -e p "$file")" ]] ; then
        { echo "$header"; cat "$file"; } > "${file}._new"
        mv "${file}._new" "${file}"
    fi
done
