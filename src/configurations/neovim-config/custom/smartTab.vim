onoremap <silent> <tab> :<c-u>call SmartTab('j')<cr>
nnoremap <silent> <tab> :<c-u>call SmartTab('j')<cr>
vnoremap <silent> <tab> :<c-u>call SmartTab_Visual('j')<cr>

onoremap <silent> <s-tab> :<c-u>call SmartTab('k')<cr>
nnoremap <silent> <s-tab> :<c-u>call SmartTab('k')<cr>
vnoremap <silent> <s-tab> :<c-u>call SmartTab_Visual('k')<cr>

" Direction is either 'j' or 'k'
func! SmartTab(direction)
  let beforeLine = line('.')
  normal! ^

  " Add to jump list
  execute "normal!".line('.')."G"

  execute "normal!".GetStopLine(a:direction)."G^"

  if beforeLine == line('.')
    execute "normal!".a:direction."^"
  endif
endfunc

func! SmartTab_Visual(direction)
  normal gv
  :call SmartTab(a:direction)
endfunc

