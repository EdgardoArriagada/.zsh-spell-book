set clipboard=unnamedplus "Share system clipboard (require xdotool)

" Paste many times over selected text
" withoud yanking it
xnoremap <expr> p 'pgv"'.v:register.'y`>'
xnoremap <expr> P 'Pgv"'.v:register.'y`>'

