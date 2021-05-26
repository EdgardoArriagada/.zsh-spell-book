onoremap <silent> <tab> :<c-u>call SmartTab('j')<cr>
nnoremap <silent> <tab> :<c-u>call SmartTab('j')<cr>
vnoremap <silent> <tab> :<c-u>call SmartTab_Visual('j')<cr>

onoremap <silent> <s-tab> :<c-u>call SmartTab('k')<cr>
nnoremap <silent> <s-tab> :<c-u>call SmartTab('k')<cr>
vnoremap <silent> <s-tab> :<c-u>call SmartTab_Visual('k')<cr>

" Direction is either 'j' or 'k'
func! SmartTab(direction)
  let beforeCol = col('.')
  let beforeLine = line('.')
  normal! ^

  "If we were not at the beggining when pressing TAB
  if col('.') != beforeCol
    return
  endif

  " Add to jump list
  execute "normal!".line('.')."G"

  :call JumpUntilNotEmptyLine(a:direction)

  if a:direction == 'j'
    :call s:goDown()
  elseif a:direction == 'k'
    :call s:goUp()
  endif

  if beforeLine == line('.')
    execute "normal!".a:direction."^"
  endif
endfunc

func! s:goDown()
  let lastMatching = GetLastMatchingIndentDown()
  if indent(lastMatching) < indent('.')
    let lastMatching -= 1
  endif
  execute "normal!".lastMatching."G^"
endfunc

func! s:goUp()
  let lastMatching = GetLastMatchingIndentUp()
  execute "normal!".lastMatching."G^"
endfunc

func! SmartTab_Visual(direction)
  normal gv
  :call SmartTab(a:direction)
endfunc

