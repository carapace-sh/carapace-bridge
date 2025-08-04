printf '%s\n' $fpath \
| xargs -I{} find {} -name '_*' 2>/dev/null \
| xargs head -n1 2>/dev/null \
| grep -v -F \
       -e ' -k ' \
       -e ' -K ' \
       -e '#autoload' \
| sed -e 's/ -[^ ]\+//g' \
      -e 's/^#compdef //' \
| tr " " "\n" \
| grep -v \
       -e '=' \
       -e '\[' \
       -e '\*' \
       -e '#' \
       -e '^_' \
       -e '^-$' \
       -e '^.$' \
| sort \
| uniq 
