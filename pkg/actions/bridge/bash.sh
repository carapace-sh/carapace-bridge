#!/usr/bin/env bash
set +o history # turn off history

COMP_WORDS=($COMP_LINE)
if [ "${COMP_LINE: -1}" = " " ]; then
  COMP_WORDS+=("")
fi
COMP_CWORD=$((${#COMP_WORDS[@]} - 1))
COMP_POINT=${#COMP_LINE}

# bash-completions
[ -f /data/data/com.termux/files/usr/share/bash-completion/bash_completion ] && source /data/data/com.termux/files/usr/share/bash-completion/bash_completion # termux
[ -f /usr/local/share/bash-completion/bash_completion ] && source /usr/local/share/bash-completion/bash_completion # osx
[ -f /usr/share/bash-completion/bash_completion ] && source /usr/share/bash-completion/bash_completion # linux

__load_completion "${COMP_WORDS[0]}"

$"$(complete -p "${COMP_WORDS[0]}" | sed -r 's/.* -F ([^ ]+).*/\1/')"

for i in "${COMPREPLY[@]}"; do
  if [[ -d "${i}" && "${i}" != */ ]]; then
    echo "${i}/"
  else
    echo "${i}"
  fi
done
