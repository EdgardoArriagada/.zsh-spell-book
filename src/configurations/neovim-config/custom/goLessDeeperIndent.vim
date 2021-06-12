onoremap <silent> <bs> :<c-u>call GoLessDeeperIndent('k')<cr>
nnoremap <silent> <bs> :<c-u>call GoLessDeeperIndent('k')<cr>
vnoremap <silent> <bs> :<c-u>call GoLessDeeperIndent_Visual('k')<cr>

onoremap <silent> <enter> :<c-u>call GoLessDeeperIndent('j')<cr>
nnoremap <silent> <enter> :<c-u>call GoLessDeeperIndent('j')<cr>
vnoremap <silent> <enter> :<c-u>call GoLessDeeperIndent_Visual('j')<cr>

func! GoLessDeeperIndent_Visual(direction)
  normal gv
  call GoLessDeeperIndent(a:direction)
endfunc

func! s:getSameIndentLine(direction)
  if a:direction == 'j'
    return GetSameIndentLineDown()
  elseif a:direction == 'k'
    return GetSameIndentLineUp()
  endif
endfunc

func! s:getLastMatchingIndent(direction)
  if a:direction == 'j'
    return GetLastMatchingIndentDown()
  elseif a:direction == 'k'
    return GetLastMatchingIndentUp()
  endif
endfunc


func! GoLessDeeperIndent(direction)
  " Go to beggin of line and add to jump list
  execute "normal!".line('.')."G^"

  :call JumpUntilNotEmptyLine(a:direction)

  let originalInent = indent('.')

  if originalInent == 0
    let lastLine = s:getSameIndentLine(a:direction)
    execute "normal!".lastLine."G^"
    return
  endif

  while col('.') > 1
    let lastMatchinLine = s:getLastMatchingIndent(a:direction)
    execute "normal!".lastMatchinLine."G^"

    if indent('.') < originalInent
      break
    endif
  endwhile
endfunc

