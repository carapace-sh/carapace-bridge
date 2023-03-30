#!/bin/bash -i
# COMP_LINE="$1"
COMP_WORDS=($COMP_LINE)
if [ "${COMP_LINE: -1}" = " " ]; then
  COMP_WORDS+=("")
fi
COMP_CWORD=$((${#COMP_WORDS[@]} - 1))
COMP_POINT=${#COMP_LINE}

$"$(complete -p "${COMP_WORDS[0]}" | sed -r 's/.* -F ([^ ]+).*/\1/')"

for i in "${COMPREPLY[@]}"; do
  if [[ -d "${i}" && "${i}" != */ ]]; then
    echo "${i}/"
  else
    echo "${i}"
  fi
done
