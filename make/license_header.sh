#!/bin/bash

# Note that we need to have extra blank lines after the license comment header
# to prevent go doc from thinking the license header is package documentation.
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
