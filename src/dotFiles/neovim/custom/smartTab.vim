vnoremap <silent> o :<c-u>call SmartTab_Visual('j')<cr>

vnoremap <silent> O :<c-u>call SmartTab_Visual('k')<cr>

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

