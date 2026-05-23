setopt extendedglob

printf '%s\n' $fpath \
| xargs -I{} find -L {} -name '_*' 2>/dev/null \
| xargs head -n1 2>/dev/null \
| grep -v -F \
       -e '==>' \
       -e ' -k ' \
       -e ' -K ' \
       -e '#autoload' \
| sed -e 's/ -[^ ]\+//g' \
      -e 's/^#compdef //' \
| tr " " "\n" \
| grep -v \
       -e '^$' \
       -e '=' \
       -e '^_' \
       -e '^-' \
       -e '^.$' \
| while IFS= read -r compdef; do
       if [[ "$compdef" == *[\[\]\*\#\(\)\|]* ]]; then
              print -rl -- ${(M)${(k)commands}:#${~compdef}}
       else
              print -r -- "$compdef"
       fi
done \
| sort \
| uniq 
