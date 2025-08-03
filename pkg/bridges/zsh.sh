printf '%s\n' $fpath \
| xargs -I{} find {} -name '_*' \
| xargs head -q -n1 \
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
