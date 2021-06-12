onoremap <silent> - :<c-u>call LookForIndentation('k')<cr>
nnoremap <silent> - :<c-u>call LookForIndentation('k')<cr>
vnoremap <silent> - :<c-u>call LookForIndentation_Visual('k')<cr>

onoremap <silent> _ :call LookForIndentation('j')<cr>
nnoremap <silent> _ :call LookForIndentation('j')<cr>
vnoremap <silent> _ :<c-u>call LookForIndentation_Visual('j')<cr>

func! LookForIndentation_Visual(direction)
  normal gv
  call LookForIndentation(a:direction)
endfunc

func! LookForIndentation(direction)
  " Go to beggin of line and add to jump list
  execute "normal!".line('.')."G^"

  :call JumpUntilNotEmptyLine(a:direction)

  if a:direction == 'j'
    let sameIndent = GetSameIndentLineDown()
  elseif a:direction == 'k'
    let sameIndent = GetSameIndentLineUp()
  endif

  execute "normal!".sameIndent."G^"
endfunc

