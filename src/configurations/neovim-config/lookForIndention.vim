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
  normal! ^m'

  " Search for a non empty to begin with
  while IsEmptyLine()
    exec 'normal! '.a:direction.'^'
  endwhile

  let l:originalValidFirstCol = col('.')
  let l:endOfFile = line('$')

  while line('.') < l:endOfFile

    exec 'normal! '.a:direction.'^'

    " [Step A] Don't use IsEmptyLine for better performance
    if len(getline(".")) == 0
      continue
    endif

    let l:currentFirstCol = col('.')

    if l:originalValidFirstCol == l:currentFirstCol
      " [Step B] Discard any false positive
      " from Step A by applying regex function this time
      if IsEmptyLine()
        continue
      endif

      return
    endif
  endwhile
endfunc

