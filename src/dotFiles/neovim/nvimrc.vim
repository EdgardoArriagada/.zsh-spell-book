set clipboard=unnamedplus "Share system clipboard (require xdotool)

" Paste many times over selected text
" withoud yanking it
vnoremap <expr> ap '$hpgv"'.v:register.'y`>'
xnoremap <expr> p 'pgv"'.v:register.'y`>'
xnoremap <expr> P 'Pgv"'.v:register.'y`>'

nnoremap <expr> k (v:count > 5 ? "m'" . v:count : "") . 'k'
nnoremap <expr> j (v:count > 5 ? "m'" . v:count : "") . 'j'

func! ZSB_ComboSelect()
  call GoLessDeeperIndent('k')
  normal! V
  call LookForIndentation('j')
endfunc

" Operate over entire function
vnoremap <silent> q :<c-u>normal! \<esc>:<c-u>call ZSB_ComboSelect()<cr>
onoremap <silent> q :<c-u>normal! :<c-u>call ZSB_ComboSelect()<cr>

func! TrimWhitespace()
  let l:save = winsaveview()
  keeppatterns %s/\s\+$//e
  call winrestview(l:save)
endfunc

augroup ZSB_NVIM_GROUP
  autocmd!
  autocmd BufWritePre * :call TrimWhitespace()
augroup END

